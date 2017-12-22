import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Observable} from 'rxjs/Observable';
import {Headers, RequestOptions}    from '@angular/http'
import 'rxjs/add/operator/map';
import { Subject } from 'rxjs/Subject';


@Injectable()
export class AuthenticationService  {
    public homeObject = new Subject<any>();
    
    constructor(private http:HttpClient){}

    /**
     * Make the Login Call in API
     * @param username 
     * @param password 
     */
    login (username:string, password:string){
        let headers = new Headers({ 'Content-Type': 'application/json' });
        let options = new RequestOptions({ headers: headers });
        return  this.http.post<any>('http://hstest.neth.eu:8080/api/login',{username:username,password:password
        })
        .map(res => res);
    }

    addDataHome(data:any){
        this.homeObject.next(data);
    }
}