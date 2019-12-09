// LitElement and html are the basic required imports
import {LitElement, html, css} from 'lit-element';

// Import our other components
import './omw-edit.js';
import './omw-report.js';

// Import 3rd party webcomponents
import './full-calendar.js';

// Create a class definition for your component and extend the LitElement base class
class OmwApp extends LitElement {
  static get properties() {
    return {
      events: { type: Array },
      outputClass: { type: String },
      outputText: { type: String },
      showHelp: { type: Boolean },
      showReport: { type: Boolean }
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
    }

    :host([hidden]) { display: none; }
    :host { display: block; }
     `;
  }

  constructor() {
    super();
    this.events = [{ // this object will be "parsed" into an Event Object
      title: 'The Title', // a property!
      start: '2018-09-01', // a property!
      end: '2018-09-02' // a property! ** see important note below about 'end' **
    }];
    this.outputClass = 'black'; // should be a CSS :host class selector
    this.outputText = '';
    this.showHelp = false;
    this.showReport = false;
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
        <input type="text" autofocus id="text-input" class="form-control text-input" name="command" @keyup="${this.handleInput}"></input>
      </div>
      <div>
        <input type="checkbox" id="helpme" class="toggle" @change="${this.toggleHelp}" ?checked=${this.showHelp}></input>
        <input type="checkbox" id="reportme" class="toggle" @change="${this.toggleReport}" ?checked=${this.showReport}></input>
      </div>
      <div class="${this.outputClass}" ?hidden=${!this.showReport}>
        <span>${this.outputText}</span>
	<full-calendar events="${this.events}"></full-calendar>
      </div>
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
    const d = new Date();
    console.log(d.toISOString(), ': Command entered = ', input);
    // clear textarea for next command
    el.value = '';

    const argv = input.split(/\s/);
    const cmd = argv.shift();
    switch(cmd) {
      case 'h':
        this.omwHello();
        break;
      case 'a':
        if (argv.length > 0) {
          this.omwAdd(argv);
        }
        else {
          this.updateOutput('Add command requires task description', 'red');
        }
        break;
      case 'r':
        this.omwReport('2019-05-27', '2019-06-03', 'json').then((report, err) => {
          if (err) {
            this.showReport = false;
            this.updateOutput(err, 'red');
            console.error('OmwReport', err)
          }
          else {
            this.showReport = true;
            this.reportData = report;
          }
        });
        break;
      case 's':
        this.omwStretch();
        break;
      case 'l':
        this.omwReport('2019-05-21', '2019-05-26', 'json').then((report, err) => {
          if (err) {
            this.showReport = false;
            this.updateOutput(err, 'red');
            console.error('OmwReport', err)
          }
          else {
            this.showReport = true;
            this.reportData = report;
          }
        })
        break;
      case 'e':
        this.omwEdit();
        break;
      case 'b':
        this.omwAdd(['break', '**']);
        break;
      case '?':
        this.showReport = false;
        this.toggleHelp();
        break;
      default:
        this.updateOutput('Invalid command - try again or ? for help', 'red');
    }
  }

  // handleInput process user input, ensure the text entered is valid
  handleInput(e) {
    if (e.key === 'Enter') {
      // Cancel the default action, if needed
      e.preventDefault();
      const el = this.shadowRoot.getElementById('text-input');
      const cmd = el.value.match(/([a-zA-Z0-9,._+:@%/-?]*) ?(\*\*\*?)*/);
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

  updateOutput(data, color) {
    this.outputClass = color;
    this.outputText = data;
  }

  getHelpText() {
    return html`
	<ul>
          <li value="cmdHello">h (hello) - start day</li>
          <li value="cmdAdd">a (add) &lt;task&gt; - add &lt;task&gt; entry with current time (use at end of task, not beginning)</li>
          <li value="cmdAddBreak">a (add) &lt;task&gt; ** - add break &lt;task&gt; entry with current time (ie: a break ***) (use at end of task, not beginning)</li>
          <li value="cmdAddIgnore">a (add) &lt;task&gt;*** - add ignored &lt;task&gt; entry with current time (ie: a commuting ***) (use at end of task, not beginning)</li>
          <li value="cmdReport">r (report) &lt;task&gt;*** - display this week\'s time report')</li>
          <li value="cmdLast">l (last) - display last week\'s time report</li>
          <li value="cmdStretch">s (stretch) &lt;task&gt;*** - stretch last task to current time')</li>
          <li value="cmdEdit">e (edit) - edit current timesheet</li>
          <li value="cmdBreak">b (break) - shortcut to add break **</li>
          <li value="cmdToggle">? (help) - toggle this help text display</li>
        </ul>`
  }

  updated(changedProperties) {
    changedProperties.forEach((oldValue, propName) => {
    });
    console.log('PROP PASSED TO UPDATED');
    console.dir(changedProperties);
  }

  omwAdd(argv) {
    async () => {
      const response = await fetch('http://127.0.0.1:31337/omw/add', {
        method: 'POST',
        body: {"args": argv},
        headers: {
          'Content-Type': 'application/json'
        }
      });
      console.log(response);
    }
  }

  omwEdit() {
    async () => {
      const response = await fetch('http://127.0.0.1:31337/omw/edit', {
        method: 'GET',
      });
      console.log(response);
    }
  }

  omwHello() {
    async () => {
      const response = await fetch('http://127.0.0.1:31337/omw/hello', {
        method: 'POST',
        body: {},
        headers: {
          'Content-Type': 'application/json'
        }
      });
      console.log(response);
    }
  }

  omwReport() {
    async () => {
      const response = await fetch('http://127.0.0.1:31337/omw/report', {
        method: 'GET',
      });
      console.log(response);
    }
  }

  omwStretch() {
    async () => {
      const response = await fetch('http://127.0.0.1:31337/omw/stretch', {
        method: 'POST',
        body: {},
        headers: {
          'Content-Type': 'application/json'
        }
      });
      console.log(response);
    }
  }
}

// Register your element to custom elements registry, pass it a tag name and your class definition
// The element name must always contain at least one dash
customElements.define('omw-app', OmwApp);
