from pydantic import BaseModel

class SimilarityResponse(BaseModel):
    similarity: float