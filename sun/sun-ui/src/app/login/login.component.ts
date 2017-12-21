import {Component} from '@angular/core';
import {Router} from '@angular/router';
import {Http, Headers} from '@angular/http';
@Component({
    moduleId:module.id,
    selector: 'login',
    templateUrl:'login.component.html',
    styleUrls: ['login.component.css']
})


export class LoginComponent{
    responseApi: string;
    constructor (
        public router: Router, public http: Http
    ){}

    login(username: string, password: string) {
        const contentHeaders = new Headers();
        contentHeaders.append('Content-Type', 'application/X-www-form-urlencoded');
        console.log(contentHeaders);
        let body = JSON.stringify({ username, password });
        this.http.post('http://hstest.neth.eu:8080/api/login', body, {headers:contentHeaders})
          .subscribe(
            response => {
              if(response.json().token){
                  localStorage.setItem('id_token', response.json().token);
                  console.log(response.json());
                   this.responseApi = response.json().account_type;
              }
            },
            error => {
              console.log(error);
              console.log("error1"+error.text().message);
              this.responseApi = error.text().message;
              console.log("error1"+error.json().message);
              this.responseApi = error.json().message;
            }
          );
    }
}