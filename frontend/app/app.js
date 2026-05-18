import ChoreoModeler from 'chor-js/lib/Modeler';
import PropertiesPanelModule from 'bpmn-js-properties-panel';

import Reporter from './lib/validator/Validator.js';
import PropertiesProviderModule from './lib/properties-provider';
import MessageMetadataPadModule from './lib/context-pad';
import ParticipantsManager from './lib/participants/ParticipantsManager.js';

import xml from './diagrams/pizzaDelivery.bpmn';
import blankXml from './diagrams/newDiagram.bpmn';

let lastFile;
let isValidating = false;
let isDirty = false;
let participantsManager;

// create and configure a chor-js instance
const modeler = new ChoreoModeler({
  container: '#canvas',
  propertiesPanel: {
    parent: '#properties-panel'
  },
  additionalModules: [
    PropertiesPanelModule,
    PropertiesProviderModule,
    MessageMetadataPadModule
  ],
  keyboard: {
    bindTo: document
  }
});

// display the given model (XML representation)
async function renderModel(newXml) {
  try {
    console.log('🎨 Import XML started...');
    const result = await modeler.importXML(newXml);
    console.log('✅ Import result:', result);
    isDirty = false;
  } catch (err) {
    console.error('❌ Error importing XML:', err);
    throw err;
  }
}

// returns the file name of the diagram currently being displayed
function diagramName() {
  if (lastFile) {
    return lastFile.name;
  }
  return 'diagram.bpmn';
}

document.addEventListener('DOMContentLoaded', () => {
  // Инициализация менеджера участников ВНУТРИ DOMContentLoaded
  participantsManager = new ParticipantsManager(modeler);
  modeler.participantsManager = participantsManager;
  
  console.log('✓ ParticipantsManager initialized');
  console.log('✓ Panel element:', document.getElementById('participants-panel'));

  // download diagram as XML
  const downloadLink = document.getElementById('js-download-diagram');
  downloadLink.addEventListener('click', async () => {
    const result = await modeler.saveXML({ format: true });
    downloadLink.href =
      'data:application/bpmn20-xml;charset=UTF-8,' +
      encodeURIComponent(result.xml);
    downloadLink.download = diagramName();
    isDirty = false;
  });

  // download diagram as SVG
  const downloadSvgLink = document.getElementById('js-download-svg');
  downloadSvgLink.addEventListener('click', async () => {
    const result = await modeler.saveSVG();
    downloadSvgLink.href =
      'data:image/svg+xml;charset=UTF-8,' +
      encodeURIComponent(result.svg);
    downloadSvgLink.download = diagramName() + '.svg';
  });

  // open file dialog
  document.getElementById('js-open-file').addEventListener('click', (e) => {
    console.log('📂 Кнопка Open нажата'); // ← Добавь это
  
    e.preventDefault();
    const input = document.getElementById('file-input');
  
    if (!input) {
      console.error('❌ #file-input не найден!');
      return;
    }
  
    console.log('📂 Input element:', input);
  
    input.value = '';
    input.click();
  });

  // toggle side panels
  const panels = Array.from(document.getElementById('panel-toggle').children);
  panels.forEach(panel => {
    panel.addEventListener('click', () => {
      const targetPanelId = panel.dataset.togglePanel;
    
      panels.forEach(p => {
        p.classList.remove('active');
        const panelEl = document.getElementById(p.dataset.togglePanel);
        if (panelEl) panelEl.classList.add('hidden');
      });

      if (targetPanelId === 'participants-panel') {
        participantsManager.show();
        panel.classList.add('active');
      } else if (targetPanelId === 'properties-panel') {
        document.getElementById('properties-panel').classList.remove('hidden');
        panel.classList.add('active');
      }
    });
  });

  // load diagram from disk
  const loadDiagram = document.getElementById('file-input');

  if (!loadDiagram) {
    console.error('❌ Элемент #file-input не найден в DOM!');
  } else {
    console.log('✅ #file-input найден, добавляем обработчик change...');
  
    // Удаляем старые обработчики (на всякий случай)
    loadDiagram.replaceWith(loadDiagram.cloneNode(true));
    const newInput = document.getElementById('file-input');
  
    newInput.addEventListener('change', async function() {
      console.log('📄 СОБЫТИЕ CHANGE СРАБОТАЛО!');
      console.log('📄 Files:', this.files);
    
      const file = this.files[0];
      if (!file) {
        console.warn('⚠️ Файл не выбран');
        return;
      }

      console.log('📄 Имя файла:', file.name);
      console.log('🔄 Чтение файла...');
    
      const reader = new FileReader();

      reader.onload = async function(e) {
        try {
          console.log('✅ Файл прочитан, длина:', e.target.result.length);
          await renderModel(e.target.result);
          console.log('✅ Модель загружена!');
          this.value = '';
        } catch (err) {
          console.error('❌ Ошибка загрузки:', err);
          alert('Ошибка: ' + err.message);
        }
      }.bind(this);
    
      reader.onerror = () => {
        console.error('❌ Ошибка чтения файла');
      };

      reader.readAsText(file);
    });
  
    console.log('✅ Обработчик change добавлен');
  }
  

  // validation
  const reporter = new Reporter(modeler);
  const validateButton = document.getElementById('js-validate');

  validateButton.addEventListener('click', () => {
    isValidating = !isValidating;

    if (isValidating) {
      reporter.validateDiagram();
      validateButton.classList.add('selected');
    } else {
      reporter.clearAll();
      validateButton.classList.remove('selected');
    }
  });

  modeler.on('commandStack.changed', () => {
    if (isValidating) reporter.validateDiagram();
    isDirty = true;
  });

  modeler.on('import.render.complete', () => {
    if (isValidating) reporter.validateDiagram();
  });
});

// expose for debugging
window.bpmnjs = modeler;

window.addEventListener('beforeunload', e => {
  if (isDirty) {
    e.preventDefault();
    e.returnValue = '';
  }
});

/**
 * =========================
 * META DATA EDITOR
 * =========================
 */

window.openMetaEditor = function(element) {
  const modal = document.getElementById('meta-modal');
  const fieldsContainer = document.getElementById('meta-fields');
  const preview = document.getElementById('meta-preview');

  let variables = [];

  const bo = element.businessObject;

  // load existing schema
  if (bo.documentation && bo.documentation[0]) {
    try {
      const parsed = JSON.parse(bo.documentation[0].text);
      variables = Object.entries(parsed.properties || {}).map(
        ([name, def]) => ({ name, type: def.type })
      );
    } catch (e) {}
  }

  function buildSchema() {
    const properties = {};
    const required = [];

    variables.forEach(v => {
      if (!v.name) return;
      properties[v.name] = { type: v.type };
      required.push(v.name);
    });

    return { properties, required };
  }

  function updatePreview() {
    preview.textContent = JSON.stringify(buildSchema(), null, 2);
  }

  function renderFields() {
    fieldsContainer.innerHTML = '';

    variables.forEach((v, idx) => {
      const row = document.createElement('div');
      row.style.display = 'flex';
      row.style.gap = '6px';
      row.style.marginBottom = '6px';

      row.innerHTML = `
        <input placeholder="name" value="${v.name || ''}" />
        <select>
          <option value="string">string</option>
          <option value="boolean">boolean</option>
          <option value="int">int</option>
        </select>
        <button>✖</button>
      `;

      const nameInput = row.querySelector('input');
      const typeSelect = row.querySelector('select');

      typeSelect.value = v.type;

      nameInput.oninput = e => {
        v.name = e.target.value;
        updatePreview();
      };

      typeSelect.onchange = e => {
        v.type = e.target.value;
        updatePreview();
      };

      row.querySelector('button').onclick = () => {
        variables.splice(idx, 1);
        renderFields();
        updatePreview();
      };

      fieldsContainer.appendChild(row);
    });
  }

  document.getElementById('meta-add-field').onclick = () => {
    variables.push({ name: '', type: 'string' });
    renderFields();
    updatePreview();
  };

  document.getElementById('meta-cancel').onclick = () => {
    modal.classList.add('hidden');
  };

  document.getElementById('meta-save').onclick = () => {
    const modeling = window.bpmnjs.get('modeling');
    const bpmnFactory = window.bpmnjs.get('bpmnFactory');

    const documentation = bpmnFactory.create('bpmn:Documentation', {
      text: JSON.stringify(buildSchema(), null, 2)
    });

    modeling.updateProperties(element, {
      documentation: [documentation]
    });

    modal.classList.add('hidden');
  };

  renderFields();
  updatePreview();
  modal.classList.remove('hidden');
};


import { ChaincodeGeneratorUI } from './lib/generator/ChaincodeGenerator.js';

let chaincodeGenerator;

document.addEventListener('DOMContentLoaded', () => {
  const generateBtn = document.getElementById('js-generate-chaincode');
  if (generateBtn) {
    chaincodeGenerator = new ChaincodeGeneratorUI(
      window.bpmnjs,
      '#msp-modal',
      '#msp-fields',
      '#msp-status'
    );
    
    generateBtn.addEventListener('click', () => {
      chaincodeGenerator.open();
    });
  }
  
  const modalGenerateBtn = document.getElementById('msp-generate');
  if (modalGenerateBtn) {
    modalGenerateBtn.addEventListener('click', () => {
      chaincodeGenerator.onGenerate();
    });
  }
  
  const modalCancelBtn = document.getElementById('msp-cancel');
  if (modalCancelBtn) {
    modalCancelBtn.addEventListener('click', () => {
      chaincodeGenerator?.close();
    });
  }
});

renderModel(blankXml); 