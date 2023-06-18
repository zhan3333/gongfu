import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CheckInTopComponent } from './check-in-top.component';

describe('CheckInTopComponent', () => {
  let component: CheckInTopComponent;
  let fixture: ComponentFixture<CheckInTopComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ CheckInTopComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(CheckInTopComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
