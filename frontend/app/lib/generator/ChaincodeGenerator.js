const API_BASE = 'http://127.0.0.1:8000';

export function getParticipantsFromModel(bpmnjs) {
  const uniqueParticipants = new Map();
  
  if (!bpmnjs) return [];
  
  try {
    const elementRegistry = bpmnjs.get('elementRegistry');
    
    elementRegistry.forEach(element => {
      const bo = element.businessObject;
      
      // Ищем все элементы типа Participant
      if (bo && (bo.$type === 'bpmn:Participant' || bo.$type === 'Participant')) {
        // Добавляем в Map (дубликаты перезапишутся по ключу id)
        if (!uniqueParticipants.has(bo.id)) {
          uniqueParticipants.set(bo.id, {
            id: bo.id,
            name: bo.name || bo.id
          });
        }
      }
    });
    
    console.log(' Найдено уникальных участников:', uniqueParticipants.size);
    
  } catch (e) {
    console.error('Error getting participants:', e);
  }
  
  // Возвращаем массив значений из Map
  return Array.from(uniqueParticipants.values());
}

export function renderMspFields(container, participants) {
  container.innerHTML = '';
  
  if (participants.length === 0) {
    container.innerHTML = `
      <div class="msp-empty-state">
        <strong>No participants found!</strong><br>
        <small>Add participants to your choreography tasks first.</small>
      </div>
    `;
    return;
  }
  
  participants.forEach((p, idx) => {
    const defaultMsp = idx % 2 === 0 ? 'Org1MSP' : 'Org2MSP';
    
    const row = document.createElement('div');
    row.className = 'msp-field-row';
    row.innerHTML = `
      <label class="msp-field-label" title="${p.id}">
        ${p.name}
      </label>
      <input 
        type="text" 
        class="msp-field-input"
        id="msp_input_${p.id}"
        data-participant="${p.id}"
        value="${defaultMsp}"
        placeholder="Org1MSP"
      />
    `;
    
    container.appendChild(row);
  });
}

export function collectMspMapping(participants) {
  const mapping = {};
  let hasError = false;
  
  participants.forEach(p => {
    const input = document.getElementById(`msp_input_${p.id}`);
    if (input) {
      const mspId = input.value.trim();
      if (!mspId) {
        input.classList.add('msp-field-error');
        hasError = true;
      } else {
        input.classList.remove('msp-field-error');
        mapping[p.id] = mspId;
      }
    }
  });
  
  return { mapping, hasError };
}

export async function generateChaincode(bpmnXml, participantMap) {
  const response = await fetch(`${API_BASE}/generate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      bpmn_xml: bpmnXml,
      participant_map: participantMap
    })
  });
  
  if (!response.ok) {
    const err = await response.json();
    throw new Error(err.detail || `HTTP ${response.status}`);
  }
  
  return await response.json();
}

export class ChaincodeGeneratorUI {
  constructor(bpmnjs, modalSelector, fieldsSelector, statusSelector) {
    this.bpmnjs = bpmnjs;
    this.modal = document.querySelector(modalSelector);
    this.fieldsContainer = document.querySelector(fieldsSelector);
    this.statusDiv = document.querySelector(statusSelector);
    this.participants = [];
  }
  
  async open() {
    this.participants = getParticipantsFromModel(this.bpmnjs);
    renderMspFields(this.fieldsContainer, this.participants);
    
    if (this.modal) this.modal.classList.remove('hidden');
    
    return this;
  }
  
  close() {
    if (this.modal) this.modal.classList.add('hidden');
    if (this.statusDiv) this.statusDiv.textContent = '';
  }
  
  async onGenerate(callback) {
    const { mapping, hasError } = collectMspMapping(this.participants);
    
    if (hasError) {
      if (this.statusDiv) {
        this.statusDiv.innerHTML = '<span class="msp-status-error">Please fill in all MSP ID fields</span>';
      }
      return;
    }
    
    if (this.participants.length === 0) {
      if (this.statusDiv) {
        this.statusDiv.innerHTML = '<span class="msp-status-error">No participants to map</span>';
      }
      return;
    }
    
    const generateBtn = document.getElementById('msp-generate');
    if (generateBtn) generateBtn.disabled = true;
    if (this.statusDiv) this.statusDiv.innerHTML = '<span class="msp-status-loading">Generating...</span>';
    
    try {
      const { xml: bpmnXml } = await this.bpmnjs.saveXML({ format: false });
      const result = await generateChaincode(bpmnXml, mapping);
      
      if (this.statusDiv) {
        this.statusDiv.innerHTML = `
          <span class="msp-status-success">
            <a href="${API_BASE}${result.download_url}" download="${result.archive_name}">
              Download ${result.archive_name}
            </a>
          </span>
        `;
      }
      
      if (typeof callback === 'function') {
        callback(null, result);
      }
      
    } catch (e) {
      if (this.statusDiv) {
        this.statusDiv.innerHTML = `<span class="msp-status-error">${e.message}</span>`;
      }
      if (typeof callback === 'function') {
        callback(e, null);
      }
    } finally {
      const generateBtn = document.getElementById('msp-generate');
      if (generateBtn) generateBtn.disabled = false;
    }
  }
}