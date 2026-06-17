from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from api import build, search, graph, health
from services.neo4j_service import neo4j_service
from config import config
import logging

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)

app = FastAPI(title="SE智图 AI Service", version="1.0.0")

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 路由
app.include_router(health.router)
app.include_router(build.router)
app.include_router(search.router)
app.include_router(graph.router)

@app.on_event("shutdown")
async def shutdown():
    neo4j_service.close()

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host=config.host, port=config.port)
