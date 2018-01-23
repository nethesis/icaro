# SunUi

This project was generated with [Angular CLI](https://github.com/angular/angular-cli) version 1.6.3.

## Must do commands
To start the sun-ui part you must follow these steps (run these commands inside sun-ui directory):
1. Run npm install
2. Run npm install -g @angular/cli

## Development server

Run `ng serve` for a dev server. Navigate to `http://localhost:4200/`. The app will automatically reload if you change any of the source files.

## Code scaffolding

Run `ng generate component component-name` to generate a new component. You can also use `ng generate directive|pipe|service|class|guard|interface|enum|module`.

## Build

Run `ng build` to build the project. The build artifacts will be stored in the `dist/` directory. Use the `-prod` flag for a production build.

## Running unit tests

Run `ng test` to execute the unit tests via [Karma](https://karma-runner.github.io).

## Running end-to-end tests

Run `ng e2e` to execute the end-to-end tests via [Protractor](http://www.protractortest.org/).

## Further help

To get more help on the Angular CLI use `ng help` or go check out the [Angular CLI README](https://github.com/angular/angular-cli/blob/master/README.md).

## Steps to implement the translation of a language

1-Add the initials of the language in the file AppComponent.ts ,
`this.translateService.addLangs(['en', 'it','initialsofthelanguage',..]);`

2-Add a file initialsofthelanguage.json in the directory /assets/i18n on project.

3-Edit the value of the language based in you translation in your json file you just added.

````bash{
  "my_english_string1": "Text to translate",
  ....
}
````