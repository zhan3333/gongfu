import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';

import { ROUTE_ANIMATIONS_ELEMENTS } from '../../../core/core.module';

import { Feature, features } from '../feature-list.data';
import { faBook } from '@fortawesome/free-solid-svg-icons';
import { TranslateModule } from '@ngx-translate/core';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { NgFor, NgClass, NgIf } from '@angular/common';
import { RtlSupportDirective } from '../../../shared/rtl-support/rtl-support.directive';

@Component({
    selector: 'anms-feature-list',
    templateUrl: './feature-list.component.html',
    styleUrls: ['./feature-list.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: true,
    imports: [RtlSupportDirective, NgFor, NgClass, MatCardModule, NgIf, MatButtonModule, FontAwesomeModule, TranslateModule]
})
export class FeatureListComponent implements OnInit {
  public faBook = faBook;
  routeAnimationsElements = ROUTE_ANIMATIONS_ELEMENTS;
  features: Feature[] = features;

  ngOnInit() {}

  openLink(link: string) {
    window.open(link, '_blank');
  }
}
