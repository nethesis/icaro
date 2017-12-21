import { HomeComponent } from './home/home.component';
import { Component } from '@angular/core';
import {LoginComponent} from './login/login.component';
import { Routes, RouterModule} from '@angular/router';
// import { NotFoundComponent } from './404/404.component'

const appRouter: Routes = [

	{path: '', redirectTo: '/login', pathMatch: 'full'},
    {path: 'login', component:LoginComponent},
    {
    	path: 'home', component:HomeComponent
    },
	
    // {path: '**', component:NotFoundComponent}
];

export const routing = RouterModule.forRoot(appRouter);