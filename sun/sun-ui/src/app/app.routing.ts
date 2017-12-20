import { HomeComponent } from './home/home.component';
import { Component } from '@angular/core';
import {LoginComponent} from './login/login.component';
import { Routes, RouterModule} from '@angular/router';


const appRouter: Routes = [
    {path: 'login', component:LoginComponent},
    {path: 'home', component:HomeComponent},
    {path: '**', redirectTo:''}
];

export const routing = RouterModule.forRoot(appRouter);