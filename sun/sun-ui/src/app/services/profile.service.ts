import { RequestOptions, Http, Headers } from '@angular/http';
import { Injectable } from '@angular/core';
import { apConfig } from '../../global';

@Injectable()
export class ProfileService {
  private apiAccountsUrl = apConfig.API_ENDPOINT_URL + '/accounts';
  private authenticationToken = localStorage.getItem('authenticationToken');

  constructor(private http: Http) {}

  getAccounts() {
    let myHeaders = new Headers();
    myHeaders.append('Content-Type', 'application/json');
    myHeaders.append('Token', this.authenticationToken);
    let myParams = new URLSearchParams();
    let options = new RequestOptions({ headers: myHeaders, params: myParams });
    return this.http
      .get(this.apiAccountsUrl, options)
      .map(res => res);
  }
}
