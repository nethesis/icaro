import { AuthenticationService } from './../_services/authentication.service';
import {Component, OnInit} from '@angular/core';
import {Router, ActivatedRoute} from '@angular/router';
import { error } from 'util';

@Component({
    moduleId:module.id,
    selector: 'router-outlet',
    templateUrl:'login.component.html'
})


export class LoginComponent{
    model:any = {};
    loading=false;
    returnUrl:string;


    constructor (
        private route: ActivatedRoute,
        private router:Router,
        private authenticationService:AuthenticationService,
    ){}


    // ngOnInit(){
    //     this.authenticationService.
    // }

    login(){
        this.loading= true;
        this.authenticationService.login(this.model.username, this.model.password)
            .subscribe(
                data=>{
                    console.log('Logged in!!');
                },
                error=>{
                    this.loading = false;
                    console.log('Fail attempt');
                }
            );
    }



}