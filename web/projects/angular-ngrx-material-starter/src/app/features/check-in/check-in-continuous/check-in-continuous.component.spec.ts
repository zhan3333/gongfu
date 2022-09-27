import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CheckInContinuousComponent } from './check-in-continuous.component';

describe('CheckInContinuousComponent', () => {
  let component: CheckInContinuousComponent;
  let fixture: ComponentFixture<CheckInContinuousComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ CheckInContinuousComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(CheckInContinuousComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
