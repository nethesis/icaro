import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Observable} from 'rxjs/Observable';
import {Headers, RequestOptions}    from '@angular/http'
import 'rxjs/add/operator/map';


@Injectable()
export class AuthenticationService  {
    constructor(private http:HttpClient){}
    login (username:string, password:string){
        let headers = new Headers({ 'Content-Type': 'application/json' });
        let options = new RequestOptions({ headers: headers });
        return  this.http.post<any>('http://hstest.neth.eu:8080/api/login',{username:username,password:password
        })
        .map(user=>{
            if(user && user.token){
                localStorage.setItem('currentUser', JSON.stringify(user));
            }
            return user;
        });
    }
}