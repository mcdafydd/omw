// LitElement and html are the basic required imports
import {LitElement, html, css} from 'lit-element';

// Import 3rd party webcomponents

// Create a class definition for your component and extend the LitElement base class
class OmwEdit extends LitElement {
  static get properties() {
    return {
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
    this.outputClass = 'black'; // should be a CSS :host class selector
  }

  // The render callback renders your element's template. This should be a pure function,
  // it should always return the same template given the same properties. It should not perform
  // any side effects such as setting properties or manipulating the DOM. See the updated
  // or first-updated examples if you need side effects.
  render() {
    // Return the template using the html template tag. This will allow lit-html to
    // interpret the dynamic parts of your template.
    return html`
      <div class="editor"></div>
    `;
  }

}

// Register your element to custom elements registry, pass it a tag name and your class definition
// The element name must always contain at least one dash
customElements.define('omw-edit', OmwEdit);
