from fastapi import HTTPException
from entities.similarity_request import SimilarityRequest
from entities.similarity_response import SimilarityResponse
from main import app, logger
from model import ProductSimilarityModel as model


@app.post("/similarity", response_model=SimilarityResponse)
async def calculate_similarity(request: SimilarityRequest):
    """
    Calculates similarity between two strings
    
    Args:
        request: Request object with two strings to compare
    
    Returns:
        SimilarityResponse: Object with comparison result and similarity coefficient
    """
    # if model.instance is None:
        # raise HTTPException(status_code=503, detail="Model not loaded")
    
    if not request.text1 or not request.text2:
        raise HTTPException(
            status_code=400,
            detail="Both texts must be non-empty"
        )
    
    try:
        # Get embeddings for both strings
        embeddings = model().encode([request.text1, request.text2])
    
        # Calculate cosine similarity
        similarity = model().similarity(embeddings[0], embeddings[1])
        
        logger.info(f"Calculated similarity: {similarity:.4f} for texts: '{request.text1[:50]}...' and '{request.text2[:50]}...'")
        
        return SimilarityResponse(
            similarity=similarity,
            text1=request.text1,
            text2=request.text2
        )
    except Exception as e:
        logger.error(f"Error calculating similarity: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error processing request: {str(e)}")