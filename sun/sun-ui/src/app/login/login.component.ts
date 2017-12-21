import {Component} from '@angular/core';
import {Router} from '@angular/router';
import {Http, Headers} from '@angular/http'
@Component({
    moduleId:module.id,
    selector: 'login',
    templateUrl:'login.component.html',
    styleUrls: ['login.component.css']
})


export class LoginComponent{
    constructor (
        public router: Router, public http: Http
    ){}

    login(username:string, password:string) {
        let contentHeaders = new Headers({'Content-Type': 'application/json'});
        let body = JSON.stringify({ username, password });
        this.http.post('http://localhost:6900/api/login',body)
          .subscribe(
            response => {
              localStorage.setItem('id_token', response.json().id_token);
              console.log(response.json());
              this.router.navigate(['home']);
            },
            error => {
             alert(error.text());
              console.log(error.text());
            }
          );
    }
}