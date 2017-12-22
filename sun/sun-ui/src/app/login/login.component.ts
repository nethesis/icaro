import { AppComponent } from './../app.component';
import { HomeComponent } from './../home/home.component';
import {Component, Input} from '@angular/core';
import {Router} from '@angular/router';
import {Http, Headers} from '@angular/http';
import { AuthenticationService } from '../_services/authentication.service';
@Component({
    moduleId:module.id,
    selector: 'login',
    templateUrl:'login.component.html',
    styleUrls: ['login.component.css']
})


export class LoginComponent{
    model: any = {};
    responseApi:string;
    constructor (
        public router: Router, public http: Http, private authenticationService: AuthenticationService
    ){}

    login() {
      this.authenticationService.login(this.model.username, this.model.password)
          .subscribe(
              data => {
                  console.log(data);
                  this.router.navigate(['/home']);
              },
              error => {
                console.log(error);
              });
  }



}