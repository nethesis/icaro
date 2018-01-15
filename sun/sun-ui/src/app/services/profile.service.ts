import { HeaderService } from './../header/header.service';
import { RequestOptions, Http, Headers } from '@angular/http';
import { Injectable } from '@angular/core';
import { apConfig } from '../../global';

@Injectable()
export class ProfileService {
  private apiAccountsUrl = apConfig.API_ENDPOINT_URL + '/accounts';
  private apiPreferencesUrl = apConfig.API_ENDPOINT_URL +
    '/preferences/accounts';
  private authenticationToken = localStorage.getItem('authenticationToken');

  constructor(private http: Http, private headerService: HeaderService) {}

  /**
   * Get Accounts list from API
   */
  getAllAccounts() {
    return this.http
      .get(this.apiAccountsUrl, this.headerService.setHeader())
      .map(res => res.json());
  }

  /**
   * Get Account in base of id
   * @param id
   */
  getAccount(id: number) {
    return this.http
      .get(this.apiAccountsUrl + '/' + id, this.headerService.setHeader())
      .map(res => res.json());
  }

  /**
   * Make the callwhich change password
   * @param account
   */
  changePassword(account: any) {
    let body = JSON.stringify(account);
    return this.http
      .put(
        this.apiAccountsUrl + '/' + '5',
        body,
        this.headerService.setHeader()
      )
      .map(res => res.json());
  }

  /**
   * Make the call which create a new account
   * @param account
   */
  postAccount(account: any) {
    let body = JSON.stringify(account);
    return this.http
      .post(this.apiAccountsUrl, body, this.headerService.setHeader())
      .map(res => res.json());
  }

  getPreferences() {
    return this.http
      .get(this.apiPreferencesUrl, this.headerService.setHeader())
      .map(res => res.json());
  }
}
