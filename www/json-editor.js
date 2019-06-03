// LitElement and html are the basic required imports
import {LitElement, html, css} from 'lit-element';

// Import 3rd party webcomponents
import {WiredInput} from 'wired-input';
import {WiredIconButton} from 'wired-icon-button';

// Other dependencies
import '@advanced-rest-client/arc-icons/arc-icons';
import '@advanced-rest-client/json-viewer/json-viewer';

// Create a class definition for your component and extend the LitElement base class
class JsonEditor extends LitElement {
  static get properties() {
    return {
      data: { type: String },
      highlightText: { type: String }
    };
  }

  static get styles() {
    return css`
    :host([hidden]) { display: none; }
    :host { display: block; }
     `;
  }

  constructor() {
    super();
    this.data = '{data: "none"}';
    this.highlightText = '';
  }

  // The render callback renders your element's template. This should be a pure function,
  // it should always return the same template given the same properties. It should not perform
  // any side effects such as setting properties or manipulating the DOM. See the updated
  // or first-updated examples if you need side effects.
  render() {
    // Return the template using the html template tag. This will allow lit-html to
    // interpret the dynamic parts of your template.
    return html`
      <h3>Content actions</h3>
        <json-viewer json="${this.data}" raw="${this.data}" query="${this.highlightText}" debug>
          <wired-icon-button slot="content-action" title="Copy content to clipboard" icon="arc:content-copy">C</wired-icon-button>
          <wired-icon-button slot="content-action" title="See raw response" icon="arc:visibility">R</wired-icon-button>
          <wired-icon-button slot="content-action" title="Save to file" icon="arc:content-copy">S</wired-icon-button>
          <wired-icon-button slot="content-action" title="(TBD) Fuzzy search Workday projects" icon="arc:content-copy">W</wired-icon-button>
          <wired-icon-button slot="content-action" title="(TBD) Get Azure Devops references" icon="arc:content-copy">G</wired-icon-button>
          <wired-input slot="content-action" type="text" id="highlight-input" class="form-control text-input" name="highlight" @keyup="${this.handleHighlight}"></wired-input>
        </json-viewer>
    `;
  }

  // handleHighlight process user input, add field content to
  // json-viewer query attribute
  handleHighlight(e) {
    // Cancel the default action, if needed
    e.preventDefault();
    var el = this.shadowRoot.getElementById('highlight-input');
    console.log('got el = ', el);
    console.dir('got event = ', e);
    this.highlightText = el.value
  }
}

// Register your element to custom elements registry, pass it a tag name and your class definition
// The element name must always contain at least one dash
customElements.define('x-json-editor', JsonEditor);
