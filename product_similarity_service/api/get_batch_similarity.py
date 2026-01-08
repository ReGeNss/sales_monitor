from fastapi import HTTPException
from entities.similarity_request import SimilarityRequest
from main import app, logger
from model import ProductSimilarityModel as model


@app.post("/batch_similarity")
async def calculate_batch_similarity(requests: list[SimilarityRequest]):
    """
    Calculates similarity for multiple string pairs simultaneously
    
    Args:
        requests: List of request objects with string pairs to compare
    
    Returns:
        List of comparison results
    """
    if model is None:
        raise HTTPException(status_code=503, detail="Model not loaded")
    
    if not requests:
        raise HTTPException(status_code=400, detail="Request list cannot be empty")
    
    try:
        results = []
        for req in requests:
            if not req.text1 or not req.text2:
                results.append({
                    "error": "Both texts must be non-empty",
                    "text1": req.text1,
                    "text2": req.text2
                })
                continue
            
            embeddings = model().encode([req.text1, req.text2])
            similarity = model().similarity(embeddings[0], embeddings[1])
            
            results.append({
                "similarity": similarity,
            })
        
        return {"results": results}
    except Exception as e:
        logger.error(f"Error calculating batch similarity: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error processing request: {str(e)}")
