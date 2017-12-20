import { Component } from '@angular/core';
import {LoginComponent} from './login/login.component';
import { Routes, RouterModule} from '@angular/router';


const appRouter: Routes = [
    {path: 'login', component:LoginComponent},


    //Otherwise direct to login 
    {path: '**', redirectTo:''}
];

export const routing = RouterModule.forRoot(appRouter);