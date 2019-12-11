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
  async handleCommand(el, input) {
    const d = new Date();
    console.log(d.toISOString(), ': Command entered = ', input);
    // clear textarea for next command
    el.value = '';

    const argv = input.split(/\s/);
    const cmd = argv.shift();
    switch(cmd) {
      case 'hello':
      case 'h':
        this.omwHello();
        break;
      case 'add':
      case 'a':
        if (argv.length > 0) {
          this.omwAdd(argv);
        }
        else {
          this.updateOutput('Add command requires task description', 'red');
        }
        break;
      case 'report':
      case 'r': {
        var err = '';
        let today = d.getFullYear() + '-' + d.getMonth()+1 + '-' + d.getDate();
        if (argv.length === 0) {
          // provide today's report
          err = await this.omwReport(today, today);
        }
        if (argv.length === 1) {
          // process argument as either:
          // YYYY-MM-DD
          // as previous day offset from current
          if (argv[0].match(/\d{1,2}/)) {
            let p = parseInt(argv[0])
            if (p === 0) {
              err = await this.omwReport(today, today);
            }
            if (p >= 0) {
              let start = today.setDate(today - p);
              err = await this.omwReport(start, today);
            }
            if (p <= 0) {
              let start = today.setDate(today + p);
              err = await this.omwReport(start, today);
            }
          }
          if (argv[0].match(/20[12][0-9]-[0-9][1-9]-[0123][1-9]/)) {
            // provided date must be today or earlier
            var start = Date.parse(argv[0]);
            if (start > Date.parse(today)) {
              console.error('Invalid date - must be today or earlier');
            }
            if (start < Date.parse(today)) {
              err = await this.omwReport(start, today);
            }
          }
        }
        if (argv.length >= 2) {
          // process arguments as either:
          // YYYY-MM-DD
          // as previous day offsets from current
          // ignore any additional trailing arguments
          if (argv[0].match(/\d{1,2}/)) {
            console.log('here');
          }
          if (argv[0].match(/20[12][0-9]-[0-9][1-9]-[0123][1-9]/)) {
            console.log('here');
          }
        }
        if (err) {
          this.showReport = false;
          this.updateOutput(err, 'red');
          console.error('OmwReport', err)
        }
        else {
          this.showReport = true;
          this.reportData = report;
        }
        break;
      }
      case 'stretch':
      case 's':
        this.omwStretch();
        break;
      case 'last':
      case 'l': {
        let day = d.getDay();
        let diff = d.getDate() - day + (day === 0 ? -6:1); // adjust when day is sunday
        let lastMonday = new Date(d.setDate(diff));
        let start = lastMonday.getFullYear() + '-' + lastMonday.getMonth()+1 + '-' + lastMonday.getDate();

        let lastFriday = lastMonday + 4;
        let end = lastFriday.getFullYear() + '-' + lastFriday.getMonth()+1 + '-' + lastFriday.getDate();

        let err = await this.omwReport(start, end);
        if (err) {
          this.showReport = false;
          this.updateOutput(err, 'red');
          console.error('OmwReport', err)
        }
        else {
          this.showReport = true;
          this.reportData = report;
        }
        break;
      }
      case 'edit':
      case 'e':
        this.omwEdit();
        break;
      case 'break':
      case 'b':
        this.omwAdd(['break', '**']);
        break;
      case 'ignore':
      case 'i':
        this.omwAdd(['ignore', '***']);
        break;
      case 'help':
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
      const cmd = el.value.match(/([a-zA-Z0-9,._+:@%\/-]+[a-zA-Z0-9,._+:@%\/\-\t ]*) ?(\*\*\*?)*/);
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
          <li value="cmdAddBreak">b (break) - shortcut to add break **</li>
          <li value="cmdAddIgnore">i (ignore) - shortcut to add ignore ***</li>
          <li value="cmdReport">r (report) &lt;task&gt;*** - display this week\'s time report')</li>
          <li value="cmdLast">l (last) - display last week\'s time report</li>
          <li value="cmdStretch">s (stretch) &lt;task&gt;*** - stretch last task to current time')</li>
          <li value="cmdEdit">e (edit) - edit current timesheet</li>
          <li value="cmdToggle">? (help) - toggle this help text display</li>
        </ul>`
  }

  updated(changedProperties) {
    console.log('PROPS PASSED TO UPDATED');
    changedProperties.forEach((oldValue, propName) => {
      console.log(`OLD = ${oldValue}`);
      console.log(`PROP NAME = ${propName}`);
    });
  }

  async omwAdd(argv) {
    await this.postApi('add', {"args": argv});
  }

  async omwEdit() {
    await this.getApi('edit');
  }

  async omwHello() {
    await this.postApi('hello', {});
  }

  // 'from' and 'to' should be YYYY-MM-DD format strings
  async omwReport(from, to) {
    await this.getApi(`report/${from}/${to}`);
  }

  async omwStretch() {
    await this.postApi('stretch', {});
  }

  async getApi(endpoint) {
    try {
      let response = await fetch(`http://localhost:31337/omw/${endpoint}`, {
        method: 'GET',
        mode: 'same-origin',
        cache: 'no-cache',
        redirect: 'error',
        referrer: 'no-referrer',
      });
      if (!response.ok) {
        this.updateOutput(`Error: Server responded with ${response.status} ${response.statusText}`, 'red');
        return {};
      }
      let tmp = await response.text();
      let data = tmp ? JSON.parse(tmp) : {};
      console.log(JSON.stringify(data));
      return data;
    } catch (error) {
      console.error('Error:', error);
    }
  }

  async postApi(endpoint, body) {
    try {
      let response = await fetch(`http://localhost:31337/omw/${endpoint}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        mode: 'same-origin',
        cache: 'no-cache',
        redirect: 'error',
        referrer: 'no-referrer',
        body: JSON.stringify(body)
      });
      if (!response.ok) {
        this.updateOutput(`Error: Server responded with ${response.status} ${response.statusText}`, 'red');
        return {};
      }
      let tmp = await response.text();
      let data = tmp ? JSON.parse(tmp) : {};
      console.log(JSON.stringify(data));
      return data;
    } catch (error) {
      console.error('Error:', error);
    }
  }
}

// Register your element to custom elements registry, pass it a tag name and your class definition
// The element name must always contain at least one dash
customElements.define('omw-app', OmwApp);
