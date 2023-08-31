import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  Inject,
  OnDestroy
} from '@angular/core';
import { Subject } from 'rxjs';
import { MatCalendar } from '@angular/material/datepicker';
import {
  DateAdapter,
  MAT_DATE_FORMATS,
  MatDateFormats
} from '@angular/material/core';
import { takeUntil } from 'rxjs/operators';
import {
  faAngleLeft,
  faAngleRight,
  faAnglesLeft,
  faAnglesRight,
  faArrowLeft,
  faArrowRight
} from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';

@Component({
    selector: 'anms-calendar-header',
    styles: [
        `
      .example-header {
        display: flex;
        align-items: center;
        padding: 0.5em;
      }
    `
    ],
    template: `
    <div class="example-header">
      <button mat-icon-button (click)="previousClicked('year')">
        <mat-icon>
          <fa-icon [icon]="faAnglesLeft"></fa-icon>
        </mat-icon>
      </button>
      <button mat-icon-button (click)="previousClicked('month')">
        <mat-icon>
          <fa-icon [icon]="faAngleLeft"></fa-icon>
        </mat-icon>
      </button>
      <span style="flex: 1;height: 1em;font-weight: 500;text-align: center;">{{
        periodLabel
      }}</span>
      <button mat-icon-button (click)="nextClicked('month')">
        <mat-icon>
          <fa-icon [icon]="faAngleRight"></fa-icon>
        </mat-icon>
      </button>
      <button mat-icon-button (click)="nextClicked('year')">
        <mat-icon>
          <fa-icon [icon]="faAnglesRight"></fa-icon>
        </mat-icon>
      </button>
    </div>
  `,
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: true,
    imports: [MatButtonModule, MatIconModule, FontAwesomeModule]
})
export class CalendarHeaderComponent<D> implements OnDestroy {
  public faArrowLeft = faArrowLeft;
  public faArrowRight = faArrowRight;
  protected readonly faAnglesRight = faAnglesRight;
  protected readonly faAnglesLeft = faAnglesLeft;
  protected readonly faAngleLeft = faAngleLeft;
  protected readonly faAngleRight = faAngleRight;
  private _destroyed = new Subject<void>();

  constructor(
    private _calendar: MatCalendar<D>,
    private _dateAdapter: DateAdapter<D>,
    @Inject(MAT_DATE_FORMATS) private _dateFormats: MatDateFormats,
    cdr: ChangeDetectorRef
  ) {
    _calendar.stateChanges
      .pipe(takeUntil(this._destroyed))
      .subscribe(() => cdr.markForCheck());
  }

  get periodLabel() {
    return this._dateAdapter
      .format(
        this._calendar.activeDate,
        this._dateFormats.display.monthYearLabel
      )
      .toLocaleUpperCase();
  }

  ngOnDestroy() {
    this._destroyed.next();
    this._destroyed.complete();
  }

  previousClicked(mode: 'month' | 'year') {
    this._calendar.activeDate =
      mode === 'month'
        ? this._dateAdapter.addCalendarMonths(this._calendar.activeDate, -1)
        : this._dateAdapter.addCalendarYears(this._calendar.activeDate, -1);
    this._calendar.monthSelected.next(this._calendar.activeDate);
  }

  nextClicked(mode: 'month' | 'year') {
    this._calendar.activeDate =
      mode === 'month'
        ? this._dateAdapter.addCalendarMonths(this._calendar.activeDate, 1)
        : this._dateAdapter.addCalendarYears(this._calendar.activeDate, 1);
    this._calendar.monthSelected.next(this._calendar.activeDate);
  }
}
