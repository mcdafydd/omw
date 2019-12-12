// LitElement and html are the basic required imports
import {LitElement, html, css} from 'lit-element';

import { Calendar } from '@fullcalendar/core';
import dayGridPlugin from '@fullcalendar/daygrid';
import interactionPlugin from '@fullcalendar/interaction';
import timeGridPlugin from '@fullcalendar/timegrid';

// Create a class definition for your component and extend the LitElement base class
class FullCalendar extends LitElement {
  static get properties() {
    return {
      calendar: { type: Object },
      header: { type: Object },
      events: { type: Array },
      weekends: { type: Boolean },
      rtl: { type: Boolean },
      hiddenDays: { type: Array },
      disableFixedWeekCount: { type: Boolean },
      weekNumbers: { type: Boolean },
      businessHours: { type: Object },
      height: { type: Number },
      defaultDate: { type: Object },
      editable: { type: Boolean },
      selectable: { type: Boolean },
      selectHelper: { type: Boolean },
      droppable: { type: Boolean },
      aspectRatio: { type: Number },
      eventStartEditable: Boolean,
      eventDurationEditable: Boolean,
      defaultView: { type: String },
      allDaySlot: { type: Boolean },
      allDayText: { type: String },
      slotDuration: { type: String },
      slotLabelInterval: { type: Object },
      snapDuration: { type: Object },
      scrollTime: { type: String },
      minTime: { type: String },
      maxTime: { type: String },
      slotEventOverlap: { type: Boolean },
      nowIndicator: { type: Boolean },
      dragRevertDuration: { type: Number },
      dragOpacity: { type: Number },
      dragScroll: { type: Boolean },
      eventOverlap: { type: Object },
      eventConstraint: { type: Object },
      locale: { type:  String },
      timezone: { type: Boolean },
      timeFormat: { type: String },
      dayNames: { type: Array },
      eventColor: { type: String },
      textColor: { type: String },
      eventRender: { type: Object },
      dayRender: { type: Object },
      options: { type: Object },
      categories: { type: Array },
      config: { type: Object },
      _scheduler: { type: Object },
    };
  }

  static get availableViews() {
    return ["dayGridMonth", "dayGridDay", "dayGridWeek", "timeGridWeek", "timeGridDay", "listYear", "listMonth", "listWeek", "listDay"];
  }

  static get defaultColors() {
    return ["#1abc9c", "#3498db", "#f1c40f", "#8e44ad", "#e74c3c", "#d35400", "#2c3e50", "#7f8c8d"];
  }

  constructor() {
    super();
    this.calendar = {};
    this.aspectRatio = 1.35;
    this.categories = [];
    this.defaultView = 'dayGridWeek';
    // Construct calendar's config object
    this._config = {
      theme: false,
      header: this.header,
      isRTL: this.rtl,
      weekends: this.weekends,
      hiddenDays: this.hiddenDays,
      fixedWeekCount: !this.disableFixedWeekCount,
      weekNumbers: this.weekNumbers,
      businessHours: this.businessHours,
      height: this.height,
      contentHeight: this.contentHeight,
      aspectRatio: this.aspectRatio,
      eventLimit: this.eventLimit,
      defaultDate: this.defaultDate,
      locale: this.locale,
      timezone: this.timezone,
      timeFormat: this.timeFormat,
      editable: this.editable,
      droppable: this.droppable,
      eventStartEditable: this.eventStartEditable,
      eventDurationEditable: this.eventDurationEditable,
      defaultView: this.defaultView,
      allDaySlot: this.allDaySlot,
      allDayText: this.allDayText,
      slotDuration: this.slotDuration,
      slotLabelInterval: this.slotLabelInterval,
      snapDuration: this.snapDuration,
      scrollTime: this.scrollTime,
      minTime: this.minTime,
      maxTime: this.maxTime,
      slotEventOverlap: this.slotEventOverlap,
      nowIndicator: this.nowIndicator,
      dragRevertDuration: this.dragRevertDuration,
      dragOpacity: this.dragOpacity,
      dragScroll: this.dragScroll,
      eventOverlap: this.eventOverlap,
      eventConstraint: this.eventConstraint,
      eventRender: this.eventRender,
      dayRender: this.dayRender,
      dayNames: this.dayNames,
      selectable: this.selectable,
      selectHelper: this.selectHelper,
      select: (start, end, jsEvent, view) => {
        this.dispatchEvent(new CustomEvent('select', { detail: { start, end, jsEvent, view } }));
      },
      dayClick: (date, jsEvent, view) => {
        // shorthand object is used to construct the detail of the dispatched event
        this.dispatchEvent(new CustomEvent('day-click', { detail: { date, jsEvent, view } }));
      },
      drop: (date, jsEvent, ui, resourceId) => {
        this.dispatchEvent(new CustomEvent('drop', { detail: { date, jsEvent, ui, resourceId } }));
      },
      eventClick: (calEvent, jsEvent, view) => {
        this.dispatchEvent(new CustomEvent('event-click', { detail: { calEvent, jsEvent, view } }));
      },
      eventMouseover: (calEvent, jsEvent, view) => {
        this.dispatchEvent(new CustomEvent('event-mouse-over', { detail: { calEvent, jsEvent, view } }));
      },
      eventMouseout: (calEvent, jsEvent, view) => {
        this.dispatchEvent(new CustomEvent('event-mouse-out', { detail: { calEvent, jsEvent, view } }));

      },
      eventDragStart: (event, jsEvent, ui, view) => {
        this.dispatchEvent(new CustomEvent('event-drag-start', { detail: { event, jsEvent, ui, view } }));

      },
      eventDragStop: (event, jsEvent, ui, view) => {
        this.dispatchEvent(new CustomEvent('event-drag-end', { detail: { event, jsEvent, ui, view } }));

      },
      eventDrop: (event, delta, revertFunc, jsEvent, ui, view) => {
        this.dispatchEvent(new CustomEvent('event-drop', { detail: { event, delta, revertFunc, jsEvent, ui, view } }));

      },
      eventResizeStart: (event, jsEvent, ui, view) => {
        this.dispatchEvent(new CustomEvent('event-resize-start', { detail: { event, jsEvent, ui, view } }));

      },
      eventResizeStop: (event, jsEvent, ui, view) => {
        this.dispatchEvent(new CustomEvent('event-resize-stop', { detail: { event, jsEvent, ui, view } }));

      },
      eventResize: (event, delta, revertFunc, jsEvent, ui, view) => {
        this.dispatchEvent(new CustomEvent('event-resize', { detail: { event, delta, revertFunc, jsEvent, ui, view } }));

      },
      viewRender: (view, element) => {
        this.dispatchEvent(new CustomEvent('view-render', { detail: { view, element } }));
        if (view.activeRange == null) return;

        const start = view.activeRange.start.toISOString();
        const end = view.activeRange.end.toISOString();
        console.log('HI DAVID', this.events);
        this.dispatchEvent(new CustomEvent('date-range-changed', { detail: { start, end } }));
      },
      viewDestroy: (view, element) => {
        this.dispatchEvent(new CustomEvent('view-destroy', { detail: { view, element } }));

      }

    }
  }

  static get styles() {
      return css`
	      :host {
          display: block;
        }

        .hide {
          display: none !important;
        }

        .show {
          display: block !important;
        }

        .fc-scroller {
          height: auto !important;
          overflow-y: hidden !important;
        }

        .breakEntry {
          color: lavender;
        }

        .ignoreEntry {
          color: beige;
        }`
  }

  render() {
    return html`
      <link href='../../node_modules/@fullcalendar/core/main.css' rel='stylesheet' />
      <link href='../../node_modules/@fullcalendar/daygrid/main.css' rel='stylesheet' />
      <link href='../../node_modules/@fullcalendar/timegrid/main.css' rel='stylesheet' />

      <div id="calendar"></div>`
  }

  connectedCallback() {
    super.connectedCallback();
  }

  firstUpdated(){
    const config={
      plugins: [dayGridPlugin, interactionPlugin, timeGridPlugin ],
      header: {left:"prev,next",center:"dayGridMonth,timeGridWeek,timeGridDay",right:"today"},
      editable: true,
      selectable: true,
      // select: info=>console.log('select'),
      // dateClick: info=>console.log('dateClick'),
      defaultView: 'dayGridWeek',
      eventSources: [
        {
          url: 'http://localhost:31337/omw/report',
          method: 'GET',
          startParam: 'start',
          endParam: 'end',
          extraParams: {
            format: 'fc'
          },
          failure: function() {
            console.error('there was an error while fetching events!');
          },
          color: 'pink',
          textColor: 'black'
        }
      ]
    };

    const calendarEl = this.shadowRoot.getElementById('calendar');
    const calendar=new Calendar(calendarEl, config);
    this.calendar = calendar;
    console.log('EVENTS = ', this.events);
   // this._updateEvents(this.events);
    calendar.render();
    console.dir(calendarEl);
  }


  /**
  * Interates over each event from the events property
  * in order to update its structure and data.
  *
  */
  _updateEvents(events) {
    console.log('UPDATE = ', this.events, ' AND ', events)
    const updatedEvents = [...events];
    updatedEvents.forEach(event => {
      if (event.category) {
        const eventCategory = this.categories.find(category => event.category.toUpperCase() === category.label.toUpperCase());
        if (eventCategory) {
          event.backgroundColor = eventCategory.color || event.backgroundColor;
          event.borderColor = eventCategory.color || event.borderColor;
          if (eventCategory.hidden === true) {
            event.className = ['hide'];
          }
        }

      }
    });
    console.log('CAL = ');
    console.dir(this.calendar);
    this.calendar.render();
    return updatedEvents;
  }

  /**
  * Adds events to the scheduler.
  *
  * */
  _refetchEvents() {
    const schedulerEvents = this._updateEvents(this.events);
    this._scheduler.fullCalendar('addEventSource', { events: schedulerEvents, color: this.eventColor, textColor: this.textColor });
  }

  /**
  * Removes events from the scheduler
  *
  * */
  _removeEvents() {
    this._scheduler.fullCalendar('removeEvents');
  }

  /**
  * Event color is calculated based on category property.
  *
  */
  _getEventColor() {

  }

  /**
  * Returns the current view of the scheduler.
  *
  * @return {string} the current view
  */
  getView() {
    return this._scheduler.fullCalendar('getView');
  }

  /**
  * Changes the scheduler view.
  *
  * @param {string} viewname the view name
  * @return {void}
  */
  changeView(viewName) {
    this._scheduler.fullCalendar('changeView', viewName);
  }

  /**
  * Get the current date of the scheduler.
  *
  * @return {string} a string representation of the current date
  */
  getDate() {
    return this._scheduler.fullCalendar('getDate');
  }

  /**
  * Navigates to the given date
  *
  * @param {string} target date
  * @return {void}
  */
  goToDate(date) {
    this._scheduler.fullCalendar('gotoDate', date);
  }

  /**
  * Navigates to the previous date based on the current view.
  *
  * @return {void}
  */
  prev() {
    this._scheduler.fullCalendar('prev');
  }

  /**
  * Navigates to the next date based on the current view.
  *
  * @return {void}
  */
  next() {
    this._scheduler.fullCalendar('next');
  }

  /**
  * Navigates to the previous year.
  *
  * @return {void}
  */
  prevYear() {
    this._scheduler.fullCalendar('prevYear');
  }

  /**
  * Navigates to the next year.
  *
  * @return {void}
  */
  nextYear() {
    this._scheduler.fullCalendar('nextYear');
  }


  /**
  * Navigates to current date
  *
  * @return {void}
  */
  today() {
    this._scheduler.fullCalendar('today');
  }

  /**
  * Increments the date displayed on the scheduler by the given duration
  *
  * @param {string} duration in milliseconds
  * @return {void}
  */
  incrementDate(duration) {
    this.schedule.fullCalendar('incrementDate', duration);
  }

  /**
  * Fired when a category is clicked
  *
  */
  _categoryClicked(e) {
    e = e.detail;
    const isCategoryDisabled = e.disabled;
    if (isCategoryDisabled) {
      this._hideEvents(e.category.label);
    }
    else {
      this._showEvents(e.category.label);
    }

  }

  /**
  * Shows Events that belongs to the given category
  *
  */
  _showEvents(categoryName) {
    const clientEvents = this._scheduler.fullCalendar('clientEvents');
    clientEvents.filter(event => event.category === categoryName)
      .forEach(event => {
        event.className.pop('hide');
        event.className.push('show');
        this._scheduler.fullCalendar('updateEvent', event);
      });
  }

  /**
  * Hides Events that belongs to the given category
  *
  */
  _hideEvents(categoryName) {
    const clientEvents = this._scheduler.fullCalendar('clientEvents');
    clientEvents.filter(event => event.category === categoryName)
      .forEach(event => {
        event.className.pop('show');
        event.className.push('hide');
        this._scheduler.fullCalendar('updateEvent', event);
      });
  }
}

// Register your element to custom elements registry, pass it a tag name and your class definition
// The element name must always contain at least one dash
customElements.define('full-calendar', FullCalendar);
