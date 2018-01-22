import { Injectable } from '@angular/core';
import { Http } from '@angular/http';
import { apConfig } from '../../global';
import { HeaderService } from './../header/header.service';

@Injectable()
export class AccountsService {
  private apiAccountsUrl = apConfig.API_ENDPOINT_URL + '/accounts';

  constructor(private http: Http, private headerService: HeaderService) { }

  /**
   * Get account by id
   * @param id 
   */
  getAccountbyId(id: number) {
    return this.http
      .get(this.apiAccountsUrl + '/' + id, this.headerService.setHeader())
      .map(res => res.json());
  }

  /**
   * Get Accounts list from API
   */
  getAllAccounts() {
    return this.http
      .get(this.apiAccountsUrl, this.headerService.setHeader())
      .map(res => res.json());
  }

  /**
   * Make api call which delete account
   * @param id 
   */
  deleteAccount(id: number) {
    return this.http
      .delete(this.apiAccountsUrl + '/' + id, this.headerService.setHeader())
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

  /**
   * Make api call which delete account
   * @param account 
   */
  editAccount(account: any) {
    let body = JSON.stringify(account);
    return this.http
      .put(
      this.apiAccountsUrl + '/' + account.id,
      body,
      this.headerService.setHeader()
      )
      .map(res => res.json());
  }
}
