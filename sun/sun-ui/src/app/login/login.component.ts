import { AppComponent } from './../app.component';
import { HomeComponent } from './../home/home.component';
import { Component, Input, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Http, Headers } from '@angular/http';
import { AuthenticationService } from '../services/authentication.service';
import { transition } from '@angular/core/src/animation/dsl';
import { NgForm } from '@angular/forms';
import { apConfig } from '../../global';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {
  model: any = {};
  messageResponseFromApi: string;
  isLogged = true;
  isUsernameWrong = true;
  isPasswordWrong = true;
  mainTitle = apConfig.MAIN_TITLE;
  logo = apConfig.LOGO_URL;

  constructor(
    public router: Router,
    public http: Http,
    private authenticationService: AuthenticationService
  ) {}

  ngOnInit() {
    this.checkIsLogged();
  }



  /**
   * Call authentication method from authenticationService,
   * make the login call in API
   */
  login(f: NgForm) {
    this.authenticationService
      .authentication(this.model.username, this.model.password)
      .subscribe(
        data => {
          this.router.navigate(['/home']);
          localStorage.setItem('account_type', data.account_type);
          // Reset form, when user is logged in!
          f.reset();
        },
        error => {
          this.messageResponseFromApi = error.error.message;
          if (this.messageResponseFromApi === 'Password is invalid') {
            this.switchAlertsMethod('password');
          } else {
            this.switchAlertsMethod('username');
          }
          this.isLogged = false;
        }
      );
  }

  /**
   * This method change the error alerts in base of API Response
   * @param wrongType
   */
  private switchAlertsMethod(wrongType: string) {
    if (wrongType === 'username') {
      this.isUsernameWrong = false;
      this.isPasswordWrong = true;
    } else {
      this.isUsernameWrong = true;
      this.isPasswordWrong = false;
    }
  }

  /**
   * Check if session isseted
   */
  private checkIsLogged() {
    if (localStorage.getItem('currentUser')) {
      this.router.navigate(['/home']);
    }
  }
}
