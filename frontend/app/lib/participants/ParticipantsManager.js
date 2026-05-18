/**
 * ParticipantsManager — управление участниками в BPMN-хореографии
 * Работает напрямую с BPMN:Choreography.participantRefs
 */
export default class ParticipantsManager {
  constructor(modeler) {
    this.modeler = modeler;
    this.elementRegistry = modeler.get('elementRegistry');
    this.modeling = modeler.get('modeling');
    this.bpmnFactory = modeler.get('bpmnFactory');

    this.panel = null;
    this.list = null;
    this.btnAdd = null;

    this._init();
  }

  _init() {
    this._findUI();
    this._bindEvents();
    this.refresh();
  }

  _findUI() {
    this.panel = document.getElementById('participants-panel');
    if (this.panel) {
      this.list = this.panel.querySelector('#participants-list');
      this.btnAdd = this.panel.querySelector('#add-participant-btn');
      if (this.btnAdd) {
        this.btnAdd.onclick = () => this._promptAdd();
      }
    }
  }

  _bindEvents() {
    // Обновляем список при любом изменении модели или импорте
    this.modeler.on(['commandStack.changed', 'import.render.complete'], () => {
      this.refresh();
    });
  }

  /**
   * Получает список УНИКАЛЬНЫХ участников из BPMN:Choreography.participantRefs
   */
  getParticipants() {
    const choreographyElements = this.elementRegistry.filter(el => el.type === 'bpmn:Choreography');
    if (choreographyElements.length === 0) return [];

    const choreographyBo = choreographyElements[0].businessObject;
    if (!choreographyBo.participantRefs) return [];

    return choreographyBo.participantRefs.map(p => ({
      id: p.id,
      name: p.name || p.id,
      bo: p // сохраняем ссылку на бизнес-объект
    }));
  }

  refresh() {
    if (!this.list) return;
    const participants = this.getParticipants();
    
    this.list.innerHTML = participants.length === 0 
      ? '<div class="empty" style="color:#999;font-size:12px;padding:10px;">Нет участников</div>' 
      : '';

    participants.forEach(p => {
      const item = document.createElement('div');
      item.className = 'participant-item';
      item.innerHTML = `
        <span class="participant-name" title="${p.name}">${p.name}</span>
        <button class="remove-participant" title="Удалить" data-id="${p.id}">✕</button>
      `;
      
      item.querySelector('.remove-participant').onclick = (e) => {
        e.stopPropagation();
        this.removeParticipant(p.id);
      };

      this.list.appendChild(item);
    });
  }

  /** Добавить участника в хореографию */
  async addParticipant(name) {
    if (!name?.trim()) return;

    const choreographyElements = this.elementRegistry.filter(el => el.type === 'bpmn:Choreography');
    if (choreographyElements.length === 0) {
      alert('Сначала создайте или откройте BPMN-хореографию!');
      return;
    }

    const choreographyBo = choreographyElements[0].businessObject;
    const currentRefs = choreographyBo.participantRefs || [];

    // Создаем новый бизнес-объект участника
    const newParticipant = this.bpmnFactory.create('bpmn:Participant', {
      id: 'Participant_' + Date.now(),
      name: name.trim(),
      processRef: this.bpmnFactory.create('bpmn:Process')
    });

    // Корректно обновляем массив participantRefs в модели
    this.modeling.updateModdleProperties(choreographyElements[0], choreographyBo, {
      participantRefs: [...currentRefs, newParticipant]
    });

    this.refresh();
  }

  /** Удалить участника из хореографии */
  removeParticipant(participantId) {
    const choreographyElements = this.elementRegistry.filter(el => el.type === 'bpmn:Choreography');
    if (choreographyElements.length === 0) return;

    const choreographyBo = choreographyElements[0].businessObject;
    const currentRefs = choreographyBo.participantRefs || [];
    const participantToRemove = currentRefs.find(p => p.id === participantId);
    if (!participantToRemove) return;

    // Проверка: используется ли участник в задачах
    const tasks = this.elementRegistry.filter(el => el.type === 'bpmn:ChoreographyTask');
    const isUsed = tasks.some(task => {
      const taskBo = task.businessObject;
      return (taskBo.participantRefs || []).includes(participantToRemove) ||
             taskBo.initiatingParticipantRef === participantToRemove;
    });

    if (isUsed) {
      alert(`Участник "${participantToRemove.name}" используется в задачах. Сначала измените задачи или удалите ссылки на него.`);
      return;
    }

    // Удаляем только из participantRefs. chor-js сам перерисует интерфейс.
    const newRefs = currentRefs.filter(p => p.id !== participantId);
    this.modeling.updateModdleProperties(choreographyElements[0], choreographyBo, {
      participantRefs: newRefs
    });

    this.refresh();
  }

  show() {
    this.panel?.classList.remove('hidden');
    this.refresh();
  }

  hide() {
    this.panel?.classList.add('hidden');
  }

  _promptAdd() {
    const name = prompt('Имя участника:');
    if (name) this.addParticipant(name);
  }
}