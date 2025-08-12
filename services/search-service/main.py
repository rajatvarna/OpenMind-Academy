import os
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field
from typing import List, Dict, Any
from elasticsearch import Elasticsearch

# --- FastAPI App Initialization ---
app = FastAPI(
    title="Search & Discovery Service",
    description="A service for indexing and searching platform content using Elasticsearch.",
    version="1.0.0",
)

# --- Elasticsearch Client Setup ---
def get_es_client():
    # In a real environment, you'd have more robust configuration.
    es_host = os.environ.get("ELASTICSEARCH_HOSTS", "http://elasticsearch:9200")
    try:
        # The client will not raise an error on creation if the host is unreachable.
        # It will only fail when an API call is made.
        client = Elasticsearch(hosts=[es_host])
        return client
    except Exception as e:
        print(f"Could not instantiate Elasticsearch client: {e}")
        return None

es_client = get_es_client()

# --- Pydantic Models ---
class IndexRequest(BaseModel):
    index_name: str = Field(..., description="The name of the index to add the document to (e.g., 'courses').")
    document_id: str
    document: Dict[str, Any]

class SearchResult(BaseModel):
    document_id: str
    score: float
    source: Dict[str, Any]

class SearchResponse(BaseModel):
    results: List[SearchResult]

# --- API Endpoints ---

@app.post("/index", status_code=201)
async def index_document(request: IndexRequest):
    """
    Indexes a document into the specified Elasticsearch index.
    """
    if es_client is None:
        raise HTTPException(status_code=503, detail="Search service is not connected to Elasticsearch.")

    print(f"Indexing document {request.document_id} into index '{request.index_name}'")

    # In a real application, you would make the following call:
    # try:
    #     response = es_client.index(
    #         index=request.index_name,
    #         id=request.document_id,
    #         document=request.document
    #     )
    #     return {"result": response['result']}
    # except Exception as e:
    //     raise HTTPException(status_code=500, detail=str(e))

    # Placeholder response
    return {"result": "created"}


@app.get("/search", response_model=SearchResponse)
async def search_documents(q: str, index_name: str = "courses"):
    """
    Searches for documents in a given index.
    """
    if es_client is None:
        raise HTTPException(status_code=503, detail="Search service is not connected to Elasticsearch.")

    print(f"Searching index '{index_name}' for query: '{q}'")

    # In a real application, you would make the following call:
    # query = {
    #     "multi_match": {
    #         "query": q,
    #         "fields": ["title", "description", "content"]
    #     }
    # }
    # try:
    #     response = es_client.search(index=index_name, query=query)
    #     results = [
    #         SearchResult(document_id=hit['_id'], score=hit['_score'], source=hit['_source'])
    #         for hit in response['hits']['hits']
    #     ]
    #     return SearchResponse(results=results)
    # except Exception as e:
    #     raise HTTPException(status_code=500, detail=str(e))

    # Placeholder response
    placeholder_results = [
        SearchResult(document_id="1", score=10.5, source={"title": "Introduction to Python", "description": "A course about Python."}),
        SearchResult(document_id="4", score=8.2, source={"title": "Graphic Design for Beginners", "description": "Learn about design and Python graphics."}),
    ]
    return SearchResponse(results=placeholder_results)

@app.get("/health")
async def health_check():
    """Health check endpoint."""
    # A real health check would ping the Elasticsearch cluster.
    # if es_client and es_client.ping():
    #     return {"status": "ok", "elasticsearch_status": "connected"}
    return {"status": "ok", "elasticsearch_status": "disconnected (simulated)"}

# To run: uvicorn main:app --host 0.0.0.0 --port 3004
