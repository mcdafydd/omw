// LitElement and html are the basic required imports
import {LitElement, html, css} from 'lit-element';

// Import 3rd party webcomponents
import {WiredInput} from 'wired-input';
import {WiredListbox} from 'wired-listbox';
import {WiredToggle} from 'wired-toggle';

import {JsonEditor} from './json-editor';

// Create a class definition for your component and extend the LitElement base class
class TimeTracker extends LitElement {
  static get properties() {
    return {
      showHelp: { type: Boolean },
      outputClass: { type: String },
      outputText: { type: String },
      reportData: { type: String }
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
    this.outputClass = 'black'; // should be a CSS :host class selector
    this.outputText = '';
    this.reportData = '{}';
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
      <div><wired-toggle id="helpme" class="toggle" @change="${this.doToggle}" ?checked=${this.showHelp}></wired-toggle></wired-toggle></div>
      <div class="${this.outputClass}">${this.outputText}</div>
      <x-json-editor data="${this.reportData}"></x-json-editor>
      <div id="helpText" ?hidden=${!this.showHelp}>${this.getHelpText()}</div>
    `;
  }

  doToggle() {
    this.showHelp = !this.showHelp;
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
        OmwHello();
        minimize();
        break;
      case 'a':
        if (argv.length > 0) {
          OmwAdd(argv);
          minimize();  
        }
        else {
          this.updateOutput('Add command requires task description', 'red');
        }
        break;
      case 'r':
        OmwReport('2019-05-27', '2019-06-03', 'json').then((report, err) => {
          if (err) {
            this.updateOutput(err, 'red');
            console.error('OmwReport', err)
          }
          else {
            this.reportData = report;
          }
        });
        minimize();
        break;
      case 's':
        OmwStretch();
        minimize();
        break;
      case 'l':
        OmwReport('2019-05-21', '2019-05-26', 'json').then((report) => {
          if (err) {
            this.updateOutput(err, 'red');
            console.error('OmwReport', err)
          }
          else {
            this.reportData = report;
          }
        })
        minimize();
        break;
      case 'e':
        OmwEdit();
        minimize();
        break;
      case 'b':
        OmwAdd(['break', '**']);
        minimize();
        break;
      case '?':
        this.doToggle();
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

  updateOutput(data, color) {
    this.outputClass = color;
    this.outputText = data;
  }

  getHelpText() {
    return html`
        <wired-listbox>
        <wired-item value="cmdHello">h (hello) - start day</wired-item>
        <wired-item value="cmdAdd">a (add) &lt;task&gt; - add &lt;task&gt; entry with current time (use at end of task, not beginning)</wired-item>
        <wired-item value="cmdAddBreak">a (add) &lt;task&gt; ** - add break &lt;task&gt; entry with current time (ie: a break ***) (use at end of task, not beginning)</wired-item>
        <wired-item value="cmdAddIgnore">a (add) &lt;task&gt;*** - add ignored &lt;task&gt; entry with current time (ie: a commuting ***) (use at end of task, not beginning)</wired-item>
        <wired-item value="cmdReport">r (report) &lt;task&gt;*** - display this week\'s time report')</wired-item>
        <wired-item value="cmdLast">l (last) - display last week\'s time report</wired-item>
        <wired-item value="cmdStretch">s (stretch) &lt;task&gt;*** - stretch last task to current time')</wired-item>
        <wired-item value="cmdEdit">e (edit) - edit current timesheet</wired-item>
        <wired-item value="cmdBreak">b (break) - shortcut to add break **</wired-item>
        <wired-item value="cmdToggle">? (help) - toggle this help text display</wired-item>`
  }
}

// Register your element to custom elements registry, pass it a tag name and your class definition
// The element name must always contain at least one dash
customElements.define('x-time-tracker', TimeTracker);
