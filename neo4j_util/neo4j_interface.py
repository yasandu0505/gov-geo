import os
import time
import logging
from neo4j import GraphDatabase
from dotenv import load_dotenv

class Neo4jInterface:
    """Interface to interact with Neo4j database"""
    
    def __init__(self):
        """Initialize the Neo4j connection"""
        # Load environment variables
        load_dotenv()
        
        # Get Neo4j connection details
        uri = os.getenv("NEO4J_URI", "neo4j://localhost:7687")
        username = os.getenv("NEO4J_USERNAME", "neo4j")
        password = os.getenv("NEO4J_PASSWORD", "")
        
        if not password:
            logging.warning("NEO4J_PASSWORD not set. Using empty password.")
        
        try:
            # Create Neo4j driver
            self.driver = GraphDatabase.driver(uri, auth=(username, password))
            
            # Verify connection
            self.driver.verify_connectivity()
            print("✅ Connected to Neo4j")
        except Exception as e:
            raise ConnectionError(f"Failed to connect to Neo4j: {e}")
    
    def close(self):
        """Close the Neo4j connection"""
        if hasattr(self, 'driver'):
            self.driver.close()
            print("Neo4j connection closed.")
    
    def execute_query(self, query, params=None):
        """Execute a Cypher query with parameters and return the result
        
        Args:
            query (str): The Cypher query to execute
            params (dict): Parameters for the query
            
        Returns:
            list: List of record dictionaries
        """
        if params is None:
            params = {}
        
        with self.driver.session() as session:
            result = session.run(query, params)
            records = result.data()  # Convert to a list of dictionaries
            return records
    
    def execute_query_with_retry(self, query, params=None, max_attempts=3):
        """Execute a Cypher query with retry mechanism
        
        Args:
            query (str): The Cypher query to execute
            params (dict): Parameters for the query
            max_attempts (int): Maximum number of retry attempts
            
        Returns:
            list: List of record dictionaries
        """
        if params is None:
            params = {}
        
        last_err = None
        
        for attempt in range(1, max_attempts + 1):
            print(f"Attempt {attempt} of {max_attempts} to execute query...")
            
            try:
                result = self.execute_query(query, params)
                print("Query executed successfully.")
                return result
            except Exception as e:
                last_err = e
                print(f"Error on attempt {attempt}: {e}")
                
                if attempt < max_attempts:
                    # Exponential backoff
                    wait_time = 2 * attempt
                    print(f"Waiting {wait_time} seconds before retrying...")
                    time.sleep(wait_time)
        
        raise Exception(f"Maximum retry attempts reached. Last error: {last_err}") 