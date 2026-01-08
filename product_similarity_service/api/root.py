from main import app
import model

@app.get("/")
async def root():
    """Root endpoint for service health check"""
    return {
        "status": "ok",
        "service": "Product Similarity Service",
        "model_loaded": model is not None
    }
