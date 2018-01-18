import { Injectable } from '@angular/core';
import { Http } from '@angular/http';
import { apConfig } from '../../global';
import { HeaderService } from './../header/header.service';

@Injectable()
export class AccountsService {
  private apiAccountsUrl = apConfig.API_ENDPOINT_URL + '/accounts';

  constructor(private http: Http, private headerService: HeaderService) { }


 
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

}
