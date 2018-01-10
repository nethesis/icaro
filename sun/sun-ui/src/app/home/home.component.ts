import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AuthenticationService } from './../services/authentication.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  public account_type = localStorage.getItem('account_type');
  constructor(
    private authenticationService: AuthenticationService,
    private router: Router
  ) {}

  ngOnInit() {}

  /**
   * This method, call logout procedure from service part.
   */
  logOut() {
    this.authenticationService.logout();
    this.router.navigate(['/login']);
  }
}
