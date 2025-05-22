#!/usr/bin/env python3

import sys
import os
import time

# Add the parent directory to sys.path
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from neo4j_util.neo4j_interface import Neo4jInterface

def create_government_node(neo4j_interface):
    """Create a government node with better error handling"""
    query = """
    MERGE (g:government {id: $id})
    SET g.name = $name
    RETURN g
    """
    params = {"id": "gov_01", "name": "Government of Sri Lanka"}
    
    max_attempts = 3
    for attempt in range(1, max_attempts + 1):
        try:
            print(f"Attempt {attempt} of {max_attempts} to create government node...")
            result = neo4j_interface.execute_query(query, params)
            print("Successfully created government node.")
            return result
        except Exception as e:
            print(f"Error on attempt {attempt}: {e}")
            if attempt < max_attempts:
                wait_time = 2 * attempt  # Increasing backoff
                print(f"Waiting {wait_time} seconds before retrying...")
                time.sleep(wait_time)
            else:
                print("Maximum retry attempts reached. Failed to create government node.")
                raise

def main():
    # Initialize Neo4j interface
    try:
        neo4j_interface = Neo4jInterface()
        
        # Create the government node
        gov_node = create_government_node(neo4j_interface)
        if gov_node:
            print("Government node details:")
            for record in gov_node:
                print(f"Node: {record['g']}")
                
    except Exception as e:
        print(f"Error in main function: {e}")
    finally:
        # Ensure we always close the connection
        try:
            neo4j_interface.close()
            print("Neo4j connection closed.")
        except Exception as e:
            print(f"Error closing Neo4j connection: {e}")

if __name__ == "__main__":
    main() 