import numpy as np
import time

# --- 1. Text Chunking ---
def chunk_text(text, chunk_size=500, overlap=50):
    """Splits text into smaller, overlapping chunks."""
    chunks = []
    start = 0
    while start < len(text):
        end = start + chunk_size
        chunks.append(text[start:end])
        start += chunk_size - overlap
    return chunks

# --- 2. Embedding Generation (Placeholder) ---
def get_embedding(text_chunk):
    """
    Placeholder for generating a vector embedding for a text chunk.
    In a real app, this would use a sentence-transformer model.
    """
    print(f"Generating embedding for chunk: '{text_chunk[:30]}...'")
    # Simulate a 768-dimension embedding
    return np.random.rand(768).tolist()

# --- 3. Vector Store (In-Memory Simulation) ---
class VectorStore:
    """
    A simple in-memory simulation of a vector database.
    """
    def __init__(self):
        self.vectors = {}
        self.documents = {}
        self.next_id = 0

    def add(self, text_chunk, embedding):
        """Adds a document and its vector to the store."""
        doc_id = str(self.next_id)
        self.vectors[doc_id] = embedding
        self.documents[doc_id] = text_chunk
        self.next_id += 1
        print(f"Stored chunk with id {doc_id}")

    def query(self, query_embedding, top_k=3):
        """
        Simulates a similarity search.
        In a real vector DB, this would perform a nearest neighbor search.
        Here, we just return a few random documents as context.
        """
        print("Performing similarity search...")
        # Get a random sample of document IDs
        num_docs = len(self.documents)
        if num_docs == 0:
            return []

        sample_size = min(top_k, num_docs)
        random_ids = np.random.choice(list(self.documents.keys()), sample_size, replace=False)

        return [self.documents[doc_id] for doc_id in random_ids]

# --- 4. LLM Call (Placeholder) ---
def ask_llm(question, context):
    """
    Placeholder for making a call to an LLM with the question and context.
    """
    print("Asking LLM for an answer...")

    # In a real app, you would format a prompt and send it to the LLM API
    prompt = f"""
    Based on the following context, please answer the user's question.
    If the context does not contain the answer, say so.

    Context:
    ---
    {" ".join(context)}
    ---

    Question: {question}

    Answer:
    """

    # Simulate API call latency
    time.sleep(2)

    # Return a hardcoded response
    return f"Based on the context provided, the answer to '{question}' is likely related to the retrieved documents. This is a simulated response from the LLM."


# Instantiate a global vector store for the service
vector_store = VectorStore()

def generate_quiz(text_content):
    """
    Placeholder for calling an LLM to generate a quiz from text.
    """
    print("Generating quiz from text...")

    # In a real app, you would format a prompt asking the LLM to create
    # a multiple-choice quiz based on the text_content.
    prompt = f"""
    Based on the following text, generate a 3-question multiple-choice quiz.
    Return the quiz as a JSON object with a "questions" array.
    Each question should have a "question", an "options" array, and a "correctAnswer" index.

    Text:
    ---
    {text_content[:1000]}...
    ---
    """

    time.sleep(3) # Simulate API call latency

    # Return a hardcoded placeholder quiz
    return {
        "questions": [
            {
                "question": "This is a sample question based on the text. What is the main topic?",
                "options": ["Option A", "Option B", "The Main Topic", "Option D"],
                "correctAnswer": 2,
            },
            {
                "question": "This is another sample question. True or False?",
                "options": ["True", "False"],
                "correctAnswer": 0,
            }
        ]
    }
