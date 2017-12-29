import { AppComponent } from './../app.component';
import { HomeComponent } from './../home/home.component';
import {Component, Input} from '@angular/core';
import {Router} from '@angular/router';
import {Http, Headers} from '@angular/http';
import { AuthenticationService } from '../_services/authentication.service';
import { loginResponse } from '../_models/loginResponse';
@Component({
    moduleId:module.id,
    selector: 'login',
    templateUrl:'login.component.html',
    styleUrls: ['login.component.css']
})


export class LoginComponent{
    home : Array<any> = [];
    isLogged:boolean = true;
    model: any = {};
    responseApi:string;
    constructor (
        public router: Router, public http: Http, private authenticationService: AuthenticationService, 
    ){}

    login() {
      this.authenticationService.login(this.model.username, this.model.password)
          .subscribe(
              data => {
                  this.isLogged = true;
                  this.router.navigate(['/home']);
              },
              error => {
                  this.isLogged = false;
                  this.responseApi = error.error.message;

              });

            this.authenticationService.homeObject.subscribe(
                data => this.home = [data, ...this.home]
            );
  }



}