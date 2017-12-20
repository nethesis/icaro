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
        
        return  this.http.post<any>('http://localhost:6900/api/login/',{username:username,password:password
        })
        .map(user=>{
            //login succesful if there is a jwk token in response
            if(user && user.token){
                localStorage.setItem('currentUser', JSON.stringify(user));
            }
            return user;
        });
    }

    // login1(book:any): Observable<any> {
    //     let headers = new Headers({ 'Content-Type': 'application/json' });
    //     return this.http.post<any>('http://localhost:6900/api/login/', book, headers)
    //     .map(user=>{
    //         //login succesful if there is a jwk token in response
    //         if(user && user.token){
    //             localStorage.setItem('currentUser', JSON.stringify(user));
    //         }
    //         return user;
    //     });
    // } 
}