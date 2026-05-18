import os
import sys
import shutil
import zipfile
import uuid
import tempfile
from datetime import datetime
from pathlib import Path
from fastapi import FastAPI, HTTPException
from fastapi.responses import FileResponse, JSONResponse
from pydantic import BaseModel
from typing import Dict, Optional, List
from fastapi.middleware.cors import CORSMiddleware

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))
from src.pipeline import BPMNTranslator

app = FastAPI(title="BPMN → Chaincode Translator")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:9013", "http://127.0.0.1:9013"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

CONTRACTS_DIR = Path(os.getenv("CONTRACTS_DIR", Path(__file__).parent / "generated_contracts"))
CONTRACTS_DIR.mkdir(parents=True, exist_ok=True)


class TranslationRequest(BaseModel):
    bpmn_xml: str
    participant_map: Optional[Dict[str, str]] = None
    contract_name: Optional[str] = None


@app.post("/generate")
async def generate(req: TranslationRequest):
    try:
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        print(req.participant_map)
        if req.contract_name:
            safe_name = "".join(c if c.isalnum() else "_" for c in req.contract_name)
            archive_name = f"contract_{safe_name}_{timestamp}.zip"
        else:
            archive_name = f"contract_{timestamp}.zip"

        archive_path = CONTRACTS_DIR / archive_name

        work_dir = Path(tempfile.mkdtemp())

        input_path = work_dir / "model.bpmn"
        with open(input_path, "w", encoding="utf-8") as f:
            f.write(req.bpmn_xml)

        template_dir = Path(__file__).parent / "contract_template"
        output_dir = work_dir / "contract"
        shutil.copytree(template_dir, output_dir)

        chaincode_path = output_dir / "chaincode" / "chaincode.go"
        translator = BPMNTranslator(str(input_path))
        translator.run(
            str(chaincode_path),
            participant_map=req.participant_map or {}
        )

        with zipfile.ZipFile(archive_path, 'w', zipfile.ZIP_DEFLATED) as zf:
            for root, dirs, files in os.walk(output_dir):
                for file in files:
                    file_path = Path(root) / file
                    arcname = file_path.relative_to(output_dir)
                    zf.write(file_path, arcname)

        metadata_path = archive_path.with_suffix('.json')
        with open(metadata_path, 'w') as f:
            import json
            json.dump({
                "timestamp": timestamp,
                "contract_name": req.contract_name,
                "participants": req.participant_map or {},
                "archive": archive_name
            }, f, indent=2)

        return JSONResponse({
            "status": "success",
            "archive_name": archive_name,
            "archive_path": str(archive_path),
            "download_url": f"/download/{archive_name}",
            "message": f"Contract saved to {archive_path}"
        })

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Translation failed: {str(e)}")


@app.get("/download/{filename}")
async def download_contract(filename: str):
    file_path = CONTRACTS_DIR / filename
    if not file_path.exists():
        raise HTTPException(status_code=404, detail="Contract not found")

    return FileResponse(
        str(file_path),
        filename=filename,
        media_type="application/zip"
    )


@app.get("/contracts")
async def list_contracts():
    contracts = []
    for file in sorted(CONTRACTS_DIR.glob("*.zip"), reverse=True):
        contracts.append({
            "name": file.name,
            "size": file.stat().st_size,
            "created": datetime.fromtimestamp(file.stat().st_mtime).isoformat(),
            "download_url": f"/download/{file.name}"
        })
    return {"contracts": contracts, "total": len(contracts)}


@app.get("/")
async def root():
    return {
        "status": "ok",
        "contracts_dir": str(CONTRACTS_DIR),
        "docs": "/docs"
    }