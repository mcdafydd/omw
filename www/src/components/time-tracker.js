// LitElement and html are the basic required imports
import {LitElement, html, css} from 'lit-element';

// Import 3rd party webcomponents
import {WiredInput} from 'wired-input';
import {WiredListbox} from 'wired-listbox';
import {WiredToggle} from 'wired-toggle';

import './grid.js';

// Create a class definition for your component and extend the LitElement base class
class TimeTracker extends LitElement {
  static get properties() {
    return {
      showHelp: { type: Boolean },
      showReport: { type: Boolean },
      outputClass: { type: String },
      outputText: { type: String }
    };
  }

  static get styles() {
    return css`
    .black {
      color: black;
    }

    .red {
      color: red;
    }

    .toggle {
      --wired-toggle-off-color:red;
      --wired-toggle-on-color:green;
    }

    :host([hidden]) { display: none; }
    :host { display: block; }
     `;
  }

  constructor() {
    super();
    this.showHelp = false;
    this.showReport = false;
    this.outputClass = 'black'; // should be a CSS :host class selector
    this.outputText = '';
    this.reportData = {};
  }

  // The render callback renders your element's template. This should be a pure function,
  // it should always return the same template given the same properties. It should not perform
  // any side effects such as setting properties or manipulating the DOM. See the updated
  // or first-updated examples if you need side effects.
  render() {
    // Return the template using the html template tag. This will allow lit-html to
    // interpret the dynamic parts of your template.
    return html`
      <div class="user-input">
        <wired-input type="text" autofocus id="text-input" class="form-control text-input" name="command" @keyup="${this.handleInput}"></wired-input>
      </div>
      <div>
        <wired-toggle id="helpme" class="toggle" @change="${this.toggleHelp}" ?checked=${this.showHelp}></wired-toggle>
        <wired-toggle id="reportme" class="toggle" @change="${this.toggleReport}" ?checked=${this.showReport}></wired-toggle>
      </div>
      <div class="${this.outputClass}">${this.outputText}</div>

      <link rel="stylesheet" href="../../node_modules/ag-grid-community/dist/styles/ag-grid.css">
      <link rel="stylesheet" href="../../node_modules/ag-grid-community/dist/styles/ag-theme-balham-dark.css">

      <div style="width: 800px;">
          <h1>Simple ag-Grid Polymer 3 Example</h1>
          <x-grid style="width: 100%; height: 350px;"
                  ?hidden=${!this.showReport}
          ></x-grid>
      <div id="helpText" ?hidden=${!this.showHelp}>${this.getHelpText()}</div>
    `;
  }
  toggleHelp() {
    this.showHelp = !this.showHelp;
  }

  toggleReport() {
    this.showReport = !this.showReport;
  }

  // handleCommand process user input and hide window after handling command without error
  handleCommand(el, input) {
    var d = new Date();
    console.log(d.toISOString(), ': Command entered = ', input);
    // clear textarea for next command
    el.value = '';

    var argv = input.split(/\s/);
    var cmd = argv.shift();
    switch(cmd) {
      case 'h':
      case 's':
      case 'b':
      case 'i':
        this.doCommand(cmd, argv, 'POST');
        break;
      case 'e':
        this.doCommand(cmd, argv, 'GET');
        break;
      case 'a':
        if (argv.length > 0) {
          this.doCommand(cmd, argv, 'POST');
        }
        else {
          this.updateOutput('Add command requires task description', 'red');
        }
        break;
      case 'r':
        d = new Date();
        var day = d.getDay();
        var monday = d.getDate() - day + (day == 0 ? -6:1); // adjust when day is sunday
        start = this.dateStr(d.setDate(monday));
        d2 = new Date();
        var day2 = d2.getDay();
        var friday = d2.getDate() - day2 + (day2 == 0 ? -6:1) + 4; // adjust when day is sunday
        end = this.dateStr(friday);
        argv = [start, end, 'json'];
        report, err = this.doCommand(cmd, argv, 'GET');
        if (err) {
          this.showReport = false;
          this.updateOutput(err, 'red');
          console.error('Report', err)
        }
        else {
          this.showReport = true;
          this.reportData = report;
        }
        break;
      case 'l':
        d = new Date();
        var day = d.getDay();
        var lastMonday = d.getDate() - day + (day == 0 ? -6:1) - 7; // adjust when day is sunday
        start = this.dateStr(d.setDate(lastMonday));
        d2 = new Date();
        var day2 = d2.getDay(),
        lastFriday = d2.getDate() - day2 + (day2 == 0 ? -6:1) -3; // adjust when day is sunday
        lastFriday = lastFriday - 7;
        end = this.dateStr(lastFriday);
        argv = [start, end, 'json'];
        report, err = this.doCommand(cmd, argv, 'GET');
        if (err) {
          this.showReport = false;
          this.updateOutput(err, 'red');
          console.error('Report', err)
        }
        else {
          this.showReport = true;
          this.reportData = report;
        }
        break;
      case '?':
        this.showReport = false;
        this.toggleHelp();
        break;
      default:
        this.updateOutput('Invalid command - try again or ? for help', 'red');
    }
  }

  dateStr(d) {
    return d.getFullYear() + '-' +
    ('0'+ (d.getMonth()+1)).slice(-2) + '-' +
    ('0'+ d.getDate()).slice(-2);
  }

  // handleInput process user input, ensure the text entered is valid
  handleInput(e) {
    if (e.key === 'Enter') {
      // Cancel the default action, if needed
      e.preventDefault();
      var el = this.shadowRoot.getElementById('text-input');
      var cmd = el.value.match(/([a-zA-Z0-9,._+:@%/-?]*) ?(\*\*\*?)*/);
      if (cmd === null) {
        this.updateOutput('Invalid command - try again or ? for help', 'red');
        el.value = '';
      }
      else {
        this.updateOutput('', 'black');
        this.handleCommand(el, cmd[0]);
      }
    }
  }

  async doCommand(cmd, argv, method) {
    await fetch('http://localhost:31337/omw/'.concat(cmd), {
      method: method,
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({args: argv})
    }).then(function (response) {
		  return response.json();
	  }).then(function (data) {
  		console.log(data);
	  });
  }

  updateOutput(data, color) {
    this.outputClass = color;
    this.outputText = data;
  }

  getHelpText() {
    return html`
        <wired-listbox>
        <wired-item value="cmdHello">h (hello) - start day</wired-item>
        <wired-item value="cmdAdd">a (add) &lt;task&gt; - add &lt;task&gt; entry with current time (use at end of task, not beginning)</wired-item>
        <wired-item value="cmdAddBreak">b (break) &lt;task&gt; ** - add break &lt;task&gt; entry with current time (ie: a break ***) (use at end of task, not beginning)</wired-item>
        <wired-item value="cmdAddIgnore">i (igore) &lt;task&gt;*** - add ignored &lt;task&gt; entry with current time (ie: a commuting ***) (use at end of task, not beginning)</wired-item>
        <wired-item value="cmdReport">r (report) &lt;task&gt;*** - display this week\'s time report')</wired-item>
        <wired-item value="cmdLast">l (last) - display last week\'s time report</wired-item>
        <wired-item value="cmdStretch">s (stretch) &lt;task&gt;*** - stretch last task to current time')</wired-item>
        <wired-item value="cmdEdit">e (edit) - edit current timesheet</wired-item>
        <wired-item value="cmdBreak">b (break) - shortcut to add break **</wired-item>
        <wired-item value="cmdToggle">? (help) - toggle this help text display</wired-item>`
  }

  updated(changedProperties) {
    changedProperties.forEach((oldValue, propName) => {
    });
    console.dir(changedProperties);
  }

}

// Register your element to custom elements registry, pass it a tag name and your class definition
// The element name must always contain at least one dash
customElements.define('x-time-tracker', TimeTracker);
