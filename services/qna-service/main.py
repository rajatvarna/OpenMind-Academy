from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import uvicorn

from rag_logic import chunk_text, get_embedding, vector_store, ask_llm

# --- FastAPI App Initialization ---
app = FastAPI(
    title="AI Q&A Service",
    description="A service for answering questions about course content using a RAG pipeline.",
    version="1.0.0",
)

# --- Pydantic Models for API Data Validation ---
class IndexRequest(BaseModel):
    document_id: str
    text_content: str

class IndexResponse(BaseModel):
    document_id: str
    chunks_indexed: int

class QueryRequest(BaseModel):
    question: str

class QueryResponse(BaseModel):
    answer: str
    context: list[str]

# --- API Endpoints ---

@app.post("/index", response_model=IndexResponse)
async def index_document(request: IndexRequest):
    """
    Endpoint to index new content.
    Receives text, chunks it, creates embeddings, and stores them.
    """

    print(f"Indexing request received for document: {request.document_id}")
    try:
        # 1. Chunk the text
        chunks = chunk_text(request.text_content)

        # 2. For each chunk, get embedding and add to vector store
        for chunk in chunks:
            embedding = get_embedding(chunk)
            vector_store.add(chunk, embedding)

        return IndexResponse(
            document_id=request.document_id,
            chunks_indexed=len(chunks)
        )
    except Exception as e:
        print(f"Error during indexing: {e}")
        raise HTTPException(status_code=500, detail="Failed to index document.")


@app.post("/query", response_model=QueryResponse)
async def query_service(request: QueryRequest):
    """
    Endpoint to ask a question.
    Receives a question, finds relevant context, and asks an LLM for an answer.
    """
    print(f"Query received: '{request.question}'")
    try:
        # 1. Get embedding for the question
        question_embedding = get_embedding(request.question)

        # 2. Query vector store for relevant context
        context_chunks = vector_store.query(question_embedding, top_k=3)

        if not context_chunks:
            raise HTTPException(status_code=404, detail="No relevant context found for this question.")

        # 3. Ask the LLM with the question and context
        answer = ask_llm(request.question, context_chunks)

        return QueryResponse(
            answer=answer,
            context=context_chunks
        )
    except HTTPException as e:
        raise e  # Re-raise HTTPException to preserve status code and detail
    except Exception as e:
        print(f"Error during query: {e}")
        raise HTTPException(status_code=500, detail="Failed to process query.")


@app.get("/health")
async def health_check():
    """Health check endpoint."""
    return {"status": "ok"}

# --- To run the server locally ---
if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=3003)
