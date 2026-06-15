from fastapi.routing import APIRouter

from .teachers import teachers_router

api_router = APIRouter(prefix="/api/v1")

api_router.include_router(teachers_router)