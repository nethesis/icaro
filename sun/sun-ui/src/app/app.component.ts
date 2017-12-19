import { Component } from '@angular/core';

@Component({
  selector: 'my-app',
  template: `<h1>Hello {{name}} {{mikel}}</h1>`,
})
export class AppComponent  { name = 'Angular'; mikel = 'mikel kanseci'}
