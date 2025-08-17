import numpy as np
import time
import os
import openai
import json

# --- OpenAI API Key ---
openai.api_key = os.getenv("OPENAI_API_KEY")
if not openai.api_key:
    print("Warning: OPENAI_API_KEY environment variable not set.")

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
    Asks the LLM a question with the given context.
    """
    print("Asking LLM for an answer...")

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

    try:
        response = openai.ChatCompletion.create(
            model="gpt-3.5-turbo",
            messages=[
                {"role": "system", "content": "You are a helpful assistant that answers questions based on the provided context."},
                {"role": "user", "content": prompt},
            ],
            temperature=0.2,
        )
        return response.choices[0].message['content'].strip()
    except Exception as e:
        print(f"Error calling OpenAI API: {e}")
        return "Sorry, I encountered an error trying to answer your question."


# Instantiate a global vector store for the service
vector_store = VectorStore()

def generate_quiz(text_content):
    """
    Generates a quiz from the given text content using an LLM.
    """
    print("Generating quiz from text...")

    prompt = f"""
    Based on the following text, generate a 3-question quiz with a mix of 'multiple-choice' and 'true-false' questions.
    Return the quiz as a JSON object with a "questions" array.
    Each question should have a "type" ('multiple-choice' or 'true-false'), a "question" string, an "options" array of strings, and a "correctAnswer" integer index.

    Text:
    ---
    {text_content[:2000]}
    ---

    JSON Quiz:
    """

    try:
        response = openai.ChatCompletion.create(
            model="gpt-3.5-turbo",
            messages=[
                {"role": "system", "content": "You are an assistant that creates educational quizzes in JSON format."},
                {"role": "user", "content": prompt},
            ],
            temperature=0.5,
        )
        quiz_json_string = response.choices[0].message['content'].strip()

        # The LLM might return the JSON wrapped in markdown, so we clean it up.
        if quiz_json_string.startswith("```json"):
            quiz_json_string = quiz_json_string[7:-4]

        return json.loads(quiz_json_string)

    except Exception as e:
        print(f"Error calling OpenAI API or parsing JSON: {e}")
        # Return a fallback quiz in case of an error
        return { "error": "Failed to generate quiz." }
