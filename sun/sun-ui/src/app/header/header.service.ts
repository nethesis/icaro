import { RequestOptions, Headers } from '@angular/http';
import { Injectable } from '@angular/core';

@Injectable()
export class HeaderService {
  private token = localStorage.getItem('authenticationToken');
  constructor() { }

  setHeader() {
    let myHeaders = new Headers();
    myHeaders.append('Content-Type', 'application/json');
    myHeaders.append('Token', this.token);
    let myParams = new URLSearchParams();
    let options = new RequestOptions({ headers: myHeaders, params: myParams });

    return options;
  }
}
