from fastapi import FastAPI
from typing import Optional
import os
import logging

from model import ProductSimilarityModel

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="Product Similarity Service",
    description="API for comparing similarity between two strings using sentence-transformers",
    version="1.0.0"
)

MODEL_PATH = os.getenv(
    "MODEL_PATH",
    os.path.join(os.path.dirname(os.path.dirname(__file__)), "product_similarity_model")
)

ProductSimilarityModel(MODEL_PATH)