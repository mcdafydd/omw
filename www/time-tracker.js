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
      showHelp: { type: Boolean }
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
      <div>
        <wired-toggle id="helpme" class="toggle" @change="${this.doToggle}"></wired-toggle></wired-toggle>
      </div>
      <span class="${this.outputClass}">${this.outputText}</span>
      <div id="helpText" ?hidden=${!this.showHelp}>${this.getHelpText()}</div>
    `;
  }

  doToggle() {
    this.showHelp = !this.showHelp;
  }

  // handleInput process user input
  handleInput(e) {
   if (e.key === 'Enter') {
    // Cancel the default action, if needed
    e.preventDefault();
    var el = this.shadowRoot.getElementById('text-input');
    var cmd = el.value.match(/([^harslebt]+)\s+([^\sA-z0-9_*-;]+)/g);
    if (cmd === null) {
      updateOutput('Invalid characters - alphanumeric only', 'red');
    }
    var d = new Date();
    console.log(d.toISOString(), ': Command entered = ', cmd);
    handleCommand(cmd);
    // clear textarea for next command
    this.value = '';
    switch(cmd) {
      case 'h':
        break;
      case 'a':
        break;
      case 'r':
        break;
      case 's':
        break;
      case 'l':
        break;
      case 'e':
        break;
      case 'b':
        break;
      case 't':
        break;
      default:
        break;
    }
  }
}

  updateOutput(data, color) {
    this.outputColor = color;
    this.outputText = data;
  }

  getHelpText() {
    return html`
        <wired-listbox>
        <wired-item value="cmdHello">hello (h) - start day</wired-item>
        <wired-item value="cmdAdd">add (a) &lt;task&gt; - add &lt;task&gt; entry with current time (use at end of task, not beginning)</wired-item>
        <wired-item value="cmdAddBreak">add (a) &lt;task&gt; ** - add break &lt;task&gt; entry with current time (ie: a break ***) (use at end of task, not beginning)</wired-item>
        <wired-item value="cmdAddIgnore">add (a) &lt;task&gt;*** - add ignored &lt;task&gt; entry with current time (ie: a commuting ***) (use at end of task, not beginning)</wired-item>
        <wired-item value="cmdReport">report (r) &lt;task&gt;*** - display this week\'s time report')</wired-item>
        <wired-item value="cmdLast">last (l) - display last week\'s time report</wired-item>
        <wired-item value="cmdStretch">stretch (s) &lt;task&gt;*** - stretch last task to current time')</wired-item>
        <wired-item value="cmdEdit">edit (e) - edit current timesheet</wired-item>
        <wired-item value="cmdBreak">break (b) - shortcut to add break **</wired-item>
        <wired-item value="cmdToggle">toggle (t) - toggle this help text display</wired-item>`
  }
}

// Register your element to custom elements registry, pass it a tag name and your class definition
// The element name must always contain at least one dash
customElements.define('x-time-tracker', TimeTracker);
