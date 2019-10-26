import {html, PolymerElement} from '@polymer/polymer/polymer-element.js';
import 'ag-grid-polymer';

import ClickableCellRenderer from './clickable-renderer';

class MyAgGrid extends PolymerElement {
  static get properties() {
    return {
      data: {
        type: Array
      },
      columns: {
        type: Array
      }
    };
  }

  static get template() {
      return html`
        <link rel="stylesheet" href="../../node_modules/ag-grid-community/dist/styles/ag-grid.css">
        <link rel="stylesheet" href="../../node_modules/ag-grid-community/dist/styles/ag-theme-balham-dark.css">

        <div style="width: 800px;">
            <h1>Simple ag-Grid Polymer 3 Example</h1>
            <ag-grid-polymer  style="width: 100%; height: 350px;"
                              class="ag-theme-balham-dark"
                              rowData="{{data}}"
                              columnDefs="{{columns}}"
                              components="{{components}}"
                              on-first-data-rendered="{{firstDataRendered}}"
                              on-row-data-changed="{{updateTable}}"
            ></ag-grid-polymer>
        </div>
  `;
  }

  constructor() {
    super();

    // passed in from other components
    this.data = [{
        timestamp: Date.now(),
        startTime: Date.now(),
        ignore: true,
        break: false,
        task: "Run report to populate grid",
        duration: "0h00m"
      }];
    this.columns = [
      { headerName: "To", field: "timestamp" },
      { headerName: "From", field: "startTime" },
      { headerName: "Ignore (Y/N)?", field: "ignore" },
      { headerName: "Break (Y/N)?", field: "break" },
      { headerName: "Task", field: "task" },
      { headerName: "Task Total Hours", field: "duration" },
      {
        headerName: "Clickable Component",
        field: "startTime",
        cellRendererFramework: 'clickable-renderer'
      }
    ];

    this.components = {
      clickableCellRenderer: ClickableCellRenderer,
    };
  }

  connectedCallback() {
    super.connectedCallback();
    console.log('CONNECT');
    console.dir(this);
    console.dir(this.data);
    console.dir(this.columns);
  }

  firstDataRendered(params) {
    params.api.sizeColumnsToFit();
    //params.rowData = this.data;
    console.log('IN FIRST');
    console.dir(params);
    console.dir(params.rowData);
    console.dir(params.columnDefs);
    OmwReport('2019-05-27', '2019-06-03', 'json').then((report, err) => {
      console.log('REPORT');
      console.log(report);
      if (err) {
        console.log(report);
        console.error('OmwReport', err);
        var obj = [{task: 'OmwReport error: ' + err}];
        params.api.setRowData(obj);
      }
      else {
        console.log(report);
        params.api.setRowData(report.entries);
      }
    });
  }

  updateTable(params) {
    console.log('IN UPDATE');
    console.dir(params);
    console.dir(params.data);
    console.dir(params.columns);
  }
}

customElements.define('x-grid', MyAgGrid);
