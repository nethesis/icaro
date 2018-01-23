import { ProfileService } from './../services/profile.service';
import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { Router } from '@angular/router';
import { AuthenticationService } from './../services/authentication.service';
import { apConfig } from '../../global';
import {TranslateService} from 'ng2-translate';
import { forEach } from '@angular/router/src/utils/collection';
import { concat } from 'rxjs/operator/concat';
import { langName } from '../../assets/i18n/lang';

@Component({
  encapsulation: ViewEncapsulation.None,
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  toggleOnClick: true;
  mainTitle = apConfig.MAIN_TITLE;
  logo = apConfig.LOGO_URL;
  public account_type = localStorage.getItem('account_type');
  private currentLanguage = localStorage.getItem('language');
  private id = Number(localStorage.getItem('id'));
  email: string;
  displayText: string;
  allLanguages: any[] = [];
  private langFullName = langName;

  constructor(
    private authenticationService: AuthenticationService,
    private profileService: ProfileService,
    private router: Router,
    private translateService: TranslateService
  ) {
    this.allLanguages = this.translateService.getLangs();
  }

  ngOnInit() {
    this.getFirstLanguage(this.currentLanguage);
    this.getAccount(this.id);
  }

  private getAccount(id: number) {
    this.profileService.getAccount(id).subscribe(
      data => {
        for (let key in data) {
          if (key === 'email') {
            this.email = data[key];
          }
        }
        console.log(this.email);
      },
      error => {
        console.log(error);
      }
    );
  }

  /**
   * This method, call logout procedure from service part.
   */
  logOut() {
    this.authenticationService.logout();
    this.router.navigate(['/login']);
  }

  /**
   * Toggle sidebar
   * Gall specific width to sidebar and main content.
   */
  toggleSidebar() {
    if (this.toggleOnClick) {
      return '185' + 'px';
    } else {
      return '75' + 'px';
    }
  }

  /**
   * Get language from user selection, and save it on session
   * @param flag
   */
  getLanguage(flag: string) {
    localStorage.setItem('language', flag);
    this.currentLanguage = localStorage.getItem('language');
    // Let translation know what languages should use
    this.translateService.use(this.currentLanguage);
    console.log(flag);
  }

  /**
   * Get language when component is beign loaded, and save it on session
   * @param flag
   */
  getFirstLanguage(lang: string) {
    this.currentLanguage = localStorage.getItem('language');
    // Let translation know what languages should use
    this.translateService.use(this.currentLanguage);
  }

  toggleSidebarText() {
    if (this.toggleOnClick) {
      return 'inline';
    } else {
      return 'none';
    }
  }
}
