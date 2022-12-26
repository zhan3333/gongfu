import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';

import { MatButtonModule } from '@angular/material/button';
import { MatMenuModule } from '@angular/material/menu';
import { MatSelectModule } from '@angular/material/select';
import { MatTabsModule } from '@angular/material/tabs';
import { MatInputModule } from '@angular/material/input';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatChipsModule } from '@angular/material/chips';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatCardModule } from '@angular/material/card';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatDividerModule } from '@angular/material/divider';
import { MatSliderModule } from '@angular/material/slider';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule, MatRippleModule } from '@angular/material/core';

import { FaIconLibrary, FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import {
  faAddressBook,
  faArrowLeft,
  faArrowRight,
  faBook, faBuilding,
  faCalendar,
  faCalendarAlt,
  faCaretDown,
  faCaretUp,
  faCheck,
  faEdit,
  faExclamationTriangle,
  faFile,
  faFilter,
  faLanguage,
  faLightbulb, faList,
  faPaintBrush,
  faPaperPlane,
  faPhone,
  faPlus,
  faSquare,
  faStream,
  faTasks,
  faTimes,
  faTrash, faUser,
  faVideo,
  faWindowMaximize
} from '@fortawesome/free-solid-svg-icons';
import { faGithub, faMediumM } from '@fortawesome/free-brands-svg-icons';

import { BigInputComponent } from './big-input/big-input/big-input.component';
import { BigInputActionComponent } from './big-input/big-input-action/big-input-action.component';
import { RtlSupportDirective } from './rtl-support/rtl-support.directive';
import { CalendarHeaderComponent } from './calendar-header/calendar-header.component';
import { FormControlPipe } from './form-control.pipe';
import { faBolt } from '@fortawesome/free-solid-svg-icons/faBolt';
import { BottomSheetComponent } from './bottom-sheet.component';
import { MatBottomSheetModule } from '@angular/material/bottom-sheet';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,

    TranslateModule,

    MatButtonModule,
    MatSelectModule,
    MatTabsModule,
    MatInputModule,
    MatProgressSpinnerModule,
    MatChipsModule,
    MatCardModule,
    MatCheckboxModule,
    MatListModule,
    MatMenuModule,
    MatIconModule,
    MatTooltipModule,
    MatSnackBarModule,
    MatSlideToggleModule,
    MatDividerModule,
    MatDatepickerModule,
    MatNativeDateModule,

    FontAwesomeModule,
    MatBottomSheetModule,
    MatRippleModule
  ],
  declarations: [
    BigInputComponent,
    BigInputActionComponent,
    RtlSupportDirective,
    CalendarHeaderComponent,
    FormControlPipe,
    BottomSheetComponent,
  ],
  exports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,

    TranslateModule,

    MatButtonModule,
    MatMenuModule,
    MatTabsModule,
    MatChipsModule,
    MatInputModule,
    MatProgressSpinnerModule,
    MatCheckboxModule,
    MatCardModule,
    MatListModule,
    MatSelectModule,
    MatIconModule,
    MatTooltipModule,
    MatSnackBarModule,
    MatSlideToggleModule,
    MatDividerModule,
    MatSliderModule,
    MatInputModule,
    MatDatepickerModule,
    MatNativeDateModule,

    FontAwesomeModule,

    BigInputComponent,
    BigInputActionComponent,
    RtlSupportDirective,
    FormControlPipe,
  ]
})
export class SharedModule {
  constructor(faIconLibrary: FaIconLibrary) {
    faIconLibrary.addIcons(
      faGithub,
      faMediumM,
      faPlus,
      faEdit,
      faTrash,
      faTimes,
      faCaretUp,
      faCaretDown,
      faExclamationTriangle,
      faFilter,
      faTasks,
      faCheck,
      faSquare,
      faLanguage,
      faPaintBrush,
      faLightbulb,
      faWindowMaximize,
      faStream,
      faBook,
      faFile,
      faVideo,
      faArrowLeft,
      faArrowRight,
      faPaperPlane,
      faTrash,
      faCalendar,
      faCalendarAlt,
      faBolt,
      faPhone,
      faUser,
      faBuilding,
      faList,
      faAddressBook
    );
  }
}
