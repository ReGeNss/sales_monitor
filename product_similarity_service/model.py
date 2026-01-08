import os
from sentence_transformers import SentenceTransformer

class ProductSimilarityModel: 
    _instance = None

    def __new__(cls, *args, **kwargs):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
        return cls._instance
    
    def __init__(self, model_path: str): 
        self._model_path = model_path
        self._model = None
        self.__load_model()

    def __load_model(self): 
        if not os.path.exists(self._model_path):
            raise FileNotFoundError(f"Model not found at path: {self._model_path}")
        try:
            self._model = SentenceTransformer(self._model_path)
        except Exception as e:
            raise Exception(f"Error loading model: {str(e)}")

    def similarity(self, text1: str, text2: str) -> float: 
        embeddings = self._model.encode([text1, text2])
        return self._model.similarity(embeddings[0], embeddings[1])