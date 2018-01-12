import { ProfileService } from './../services/profile.service';
import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { CardConfig, InfoStatusCardConfig } from 'patternfly-ng';

@Component({
  encapsulation: ViewEncapsulation.None,
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit {
  config: CardConfig;
  public account_type = localStorage.getItem('account_type');

  constructor(private profileService: ProfileService) {}

  ngOnInit() {
    this.profileService.getAccounts();
    this.config = {
      noPadding: false,
      topBorder: true
    } as CardConfig;
  }

  private getAllAccounts() {
    this.profileService.getAccounts().subscribe(
      data => {
        console.log(data);
      },
      error => {
        console.log(error);
      }
    );
  }
}
