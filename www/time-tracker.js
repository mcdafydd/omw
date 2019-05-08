// LitElement and html are the basic required imports
import {LitElement, html, css} from 'lit-element';

// Import 3rd party webcompoents
import {WiredInput} from 'wired-input';
import {WiredListbox} from 'wired-listbox';
import {WiredToggle} from 'wired-toggle';

// Create a class definition for your component and extend the LitElement base class
class TimeTracker extends LitElement {
  static get properties() {
    return {
      showHelp: { type: Boolean },
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
    this.outputClass = 'black'; // should be a CSS :host class selector
    this.outputText = '';
  }

  // The render callback renders your element's template. This should be a pure function,
  // it should always return the same template given the same properties. It should not perform
  // any side effects such as setting properties or manipulating the DOM. See the updated
  // or first-updated examples if you need side effects.
  render() {
    // Return the template using the html template tag. This will allow lit-html to
    // interpret the dynamic parts of your template.
    return html`
      <div><wired-input type="text" autofocus id="text-input" class="form-control text-input" name="command" @keyup="${this.handleInput}"></wired-input>
      <div><wired-toggle id="helpme" class="toggle" @change="${this.doToggle}" ?checked=${this.showHelp}></wired-toggle></wired-toggle></div>
      <div class="${this.outputClass}">${this.outputText}</div>
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
        runUtt(['hello']);
        minimize();
        break;
      case 'a':
        if (len(argv) > 0) {
          runUtt(['add'].concat(argv));
          minimize();  
        }
        else {
          this.updateOutput('Add command requires task description', 'red');
        }
        break;
      case 'r':
        runUtt(['report', '--from', 'monday', '--to', 'friday']);
        minimize();
        break;
      case 's':
        runUtt(['stretch']);
        minimize();
        break;
      case 'l':
        runUtt(['report', '--from', 'monday', '--to', 'friday']);
        minimize();
        break;
      case 'e':
        runUtt(['edit']);
        minimize();
        break;
      case 'b':
        runUtt(['add', 'break', '**']);
        minimize();
        break;
      case 't':
        this.doToggle();
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
      var cmd = el.value.match(/^([A-z?]+)(\s+[A-z0-9_*-;\s]+)?$/g);    
      if (cmd === null) {
        this.updateOutput('Invalid characters - alphanumeric only', 'red');
        el.value = '';
      }
      else {
        this.handleCommand(el, cmd[0]);  
        // Clear any previous output
        this.updateOutput('', 'black');
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
        <wired-item value="cmdToggle">t (toggle), ?) - toggle this help text display</wired-item>`
  }
}

// Register your element to custom elements registry, pass it a tag name and your class definition
// The element name must always contain at least one dash
customElements.define('x-time-tracker', TimeTracker);
